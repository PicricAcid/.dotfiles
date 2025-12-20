local comment_bp = require('comment_bp_module')

vim.keymap.set('n', '<Leader>b', function()
  comment_bp.add_rating('b')
end, { silent = true, desc = 'Add positive rating to comment' })

vim.keymap.set('n', '<Leader>q', function()
  comment_bp.add_rating('p')
end, { silent = true, desc = 'Add negative rating to comment' })

vim.api.nvim_create_user_command('CommentBPList', function()
  comment_bp.show_rating_list()
end, { desc = 'Show comment rating list in QuickFix' })

vim.api.nvim_create_user_command('CommentBPShow', function()
  comment_bp.show_popup()
end, { desc = 'Show rating popup for current comment' })

vim.api.nvim_create_user_command('CommentBPDebug', function()
  comment_bp.debug_popup()
end, { desc = 'Show debug information for current line' })

local conceal_matches = {}

local function setup_conceal()
  local bufnr = vim.api.nvim_get_current_buf()
  local ft = vim.bo.filetype

  local supported = vim.tbl_contains(
    {'python', 'c', 'cpp', 'vim', 'lua', 'sh', 'bash', 'ruby', 'java', 'javascript', 'typescript'},
    ft
  )

  if not supported then
    return
  end

  if conceal_matches[bufnr] then
    for _, match_id in ipairs(conceal_matches[bufnr]) do
      pcall(vim.fn.matchdelete, match_id)
    end
  end
  conceal_matches[bufnr] = {}

  -- Conceal compact format: b:N, p:N
  local match_id = vim.fn.matchadd('Conceal', [[\s\+b:\d\+,\s*p:\d\+]], 10, -1, {conceal = ''})
  table.insert(conceal_matches[bufnr], match_id)
 
  if vim.wo.conceallevel == 0 then
    vim.wo.conceallevel = 2
  end

  if vim.wo.concealcursor == '' then
    vim.wo.concealcursor = 'nc'
  end
end

local conceal_group = vim.api.nvim_create_augroup('CommentBPConceal', { clear = true })

vim.api.nvim_create_autocmd({'FileType', 'BufEnter', 'BufWinEnter'}, {
  group = conceal_group,
  pattern = {'*.py', '*.c', '*.cpp', '*.h', '*.hpp', '*.vim', '*.lua', '*.sh', '*.bash', '*.rb', '*.java', '*.js', '*.ts'},
  callback = function()
    vim.schedule(setup_conceal)
  end,
  desc = 'Setup conceal for comment ratings'
})

local timer = nil
local debounce_ms = 100

if vim.g.comment_bp_auto_popup == nil then
  vim.g.comment_bp_auto_popup = true
end

if vim.g.comment_bp_auto_popup then
  vim.api.nvim_create_autocmd({'CursorMoved', 'CursorMovedI'}, {
    group = vim.api.nvim_create_augroup('CommentBP', { clear = true }),
    callback = function()
      if timer then
        vim.fn.timer_stop(timer)
        timer = nil
      end

      timer = vim.fn.timer_start(debounce_ms, function()
        comment_bp.show_popup()
        timer = nil
      end)
    end,
    desc = 'Auto show rating popup on cursor move'
  })
end
