local M = {}

local popup_info = require('utils.popup_info')

local function supports_emoji()
  if vim.g.comment_bp_use_emoji ~= nil then
    return vim.g.comment_bp_use_emoji
  end

  return vim.fn.has('multi_byte') == 1
end

local function format_rating(b, p)
  if supports_emoji() then
    return string.format('👍 %d 👎 %d', b, p)
  else
    return string.format('+%d -%d', b, p)
  end
end

local function get_rating_patterns()
  return {
    compact = '%s+b:(%d+),%s*p:(%d+)',
    emoji = '%s+(👍%s*(%d+)%s*👎%s*(%d+))',
    simple = '%s+%+(%d+)%s*%-(%d+)',
  }
end

local function get_comment_patterns()
  local ft = vim.bo.filetype
  local patterns = {}

  if ft == 'c' or ft == 'cpp' or ft == 'java' or ft == 'javascript' or ft == 'typescript' then
    patterns = {
      '^%s*//(.*)$',
      '^%s*/%*(.*)%*/$'
    }
  elseif ft == 'python' or ft == 'sh' or ft == 'bash' or ft == 'ruby' then
    patterns = { '^%s*#(.*)$' }
  elseif ft == 'vim' then
    patterns = { '^%s*"(.*)$' }
  elseif ft == 'lua' then
    patterns = { '^%s*%-%-(.*)$' }
  else
    patterns = {
      '^%s*//(.*)$',
      '^%s*/%*(.*)%*/$',
      '^%s*#(.*)$',
      '^%s*"(.*)$',
      '^%s*%-%-(.*)$'
    }
  end

  return patterns
end

local function is_comment_line(line)
  local patterns = get_comment_patterns()

  for _, pattern in ipairs(patterns) do
    if line:match(pattern) then
      return true
    end
  end

  return false
end

local function extract_rating(line)
  local result = { found = false, b = 0, p = 0 }
  local patterns = get_rating_patterns()

  -- Check compact format first (b:N, p:N)
  local b_count, p_count = line:match(patterns.compact)
  if b_count and p_count then
    local b_num = tonumber(b_count)
    local p_num = tonumber(p_count)
    if b_num and p_num then
      result.found = true
      result.b = b_num
      result.p = p_num
      return result
    end
  end

  -- Check emoji format for backward compatibility
  b_count, p_count = line:match(patterns.emoji)
  if b_count and p_count then
    local b_num = tonumber(b_count)
    local p_num = tonumber(p_count)
    if b_num and p_num then
      result.found = true
      result.b = b_num
      result.p = p_num
      return result
    end
  end

  -- Check simple format for backward compatibility
  b_count, p_count = line:match(patterns.simple)
  if b_count and p_count then
    local b_num = tonumber(b_count)
    local p_num = tonumber(p_count)
    if b_num and p_num then
      result.found = true
      result.b = b_num
      result.p = p_num
      return result
    end
  end

  return result
end

local function update_rating(line, rating_type)
  local rating = extract_rating(line)
  local patterns = get_rating_patterns()

  if rating_type == 'b' then
    rating.b = rating.b + 1
  elseif rating_type == 'p' then
    rating.p = rating.p + 1
  else
    vim.api.nvim_err_writeln('Invalid rating type: ' .. rating_type)
    return line
  end

  -- Always use compact format for embedding
  local new_rating = string.format(' b:%d, p:%d', rating.b, rating.p)

  if rating.found then
    local replaced = false

    -- Replace compact format first
    if line:match(patterns.compact) then
      line = line:gsub(patterns.compact, new_rating)
      replaced = true
    end

    -- Replace old emoji format (migrate to new format)
    if not replaced and line:match(patterns.emoji) then
      line = line:gsub(patterns.emoji, new_rating)
      replaced = true
    end

    -- Replace old simple format (migrate to new format)
    if not replaced and line:match(patterns.simple) then
      line = line:gsub(patterns.simple, new_rating)
      replaced = true
    end

  else
    line = line:gsub('%s*$', '')

    if line:match('%*/$') then
      line = line:gsub('%s*%*/$', new_rating .. ' */')
    else
      line = line .. new_rating
    end
  end

  return line
end

local function close_float()
  popup_info.close()
end

function M.add_rating(rating_type)
  local line = vim.api.nvim_get_current_line()

  if not is_comment_line(line) then
    print('Not a comment line')
    return
  end

  local new_line = update_rating(line, rating_type)

  local row = vim.api.nvim_win_get_cursor(0)[1]
  vim.api.nvim_buf_set_lines(0, row - 1, row, false, { new_line })

  if rating_type == 'b' then
    if supports_emoji() then
      print('Added positive rating 👍')
    else
      print('Added positive rating (+)')
    end
  elseif rating_type == 'p' then
    if supports_emoji() then
      print('Added negative rating 👎')
    else
      print('Added negative rating (-)')
    end
  else
    vim.api.nvim_err_writeln('Invalid rating type: ' .. rating_type)
  end
end

function M.show_rating_list()
  local qf_list = {}
  local lines = vim.api.nvim_buf_get_lines(0, 0, -1, false)

  for line_num, line in ipairs(lines) do
    if is_comment_line(line) then
      local rating = extract_rating(line)
      local text

      if rating.found then
        text = string.format('%s | %s', format_rating(rating.b, rating.p), line)
      else
        text = string.format('No rating | %s', line)
      end

      table.insert(qf_list, {
        lnum = line_num,
        text = text
      })
    end
  end

  vim.fn.setqflist(qf_list)

  if #qf_list == 0 then
    print('No comments found')
  else
    vim.cmd('copen')
    print(string.format('Found %d comments', #qf_list))
  end
end

function M.show_popup()
  local line = vim.api.nvim_get_current_line()

  if not is_comment_line(line) then
    close_float()
    return
  end

  local rating = extract_rating(line)
  if not rating.found then
    close_float()
    return
  end

  local text = format_rating(rating.b, rating.p)

  local ui = vim.api.nvim_list_uis()[1]
  if not ui then
    return
  end

  local cursor_pos = vim.api.nvim_win_get_cursor(0)
  local row = cursor_pos[1]

  local width = #text + 2
  local col = ui.width - width - 1

  popup_info.popup_info(text, {
    relative = 'editor',
    row = row,
    col = col,
    width = width,
    height = 1,
    border = 'rounded',
    title = 'Rating',
    zindex = 200,
    timeout = nil  -- No timeout, show permanently
  })
end

return M
