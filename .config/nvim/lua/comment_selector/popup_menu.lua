local api = vim.api
local core = require("comment_selector.core")

local function popup_menu(popup_table, opt, callback)
  local win = 0
  local buffer = api.nvim_create_buf(false, true)
  opt.title = { { opt.title, 'DarksoulsFloatTitle' } }

  local start_index = opt.start_index
  local cursor_pos = opt.start_index
  local height = math.min(opt.height, #popup_table)
    print(height)

  win = api.nvim_open_win(buffer, true, {
    relative = opt.relative,
    row = opt.row,
    col = opt.col,
    width = opt.width,
    height = height,
    focusable = true,
    border = opt.border,
    title = opt.title,
    title_pos = 'left',
    noautocmd = true,
    zindex = opt.zindex,
    win = opt.win,
  })

  api.nvim_win_set_option(win, 'number', false)
  api.nvim_win_set_option(win, 'relativenumber', false)
  api.nvim_win_set_option(win, 'wrap', false)
  api.nvim_win_set_option(win, 'cursorline', false)
  api.nvim_win_set_option(win, 'winhighlight', 'FloatBorder:DarksoulsFloatBorder,NormalFloat:DarksoulsFloatText')

  local ns_id = api.nvim_create_namespace("popup_menu_ns")
  local current_extmark = nil

  local function window_update(pos)
    if current_extmark then
      api.nvim_buf_del_extmark(buffer, ns_id, current_extmark)
    end
    current_extmark = api.nvim_buf_set_extmark(buffer, ns_id, pos - 1, 0, {
      end_row = pos,
      hl_group = "DarksoulsFloatTextSelected",
      priority = 100,
      hl_eol = true,
    })
  end

  local function render()
    local end_index = math.min(start_index + height - 1, #popup_table)
    local display = {}
    for i = start_index, end_index do
      table.insert(display, popup_table[i])
    end
    api.nvim_buf_set_lines(buffer, 0, -1, true, display)
    window_update(cursor_pos)
  end

  local function move_cursor(dir)
    if dir == "down" then
      if cursor_pos < height and (start_index + cursor_pos - 1) < #popup_table then
        cursor_pos = cursor_pos + 1
      elseif (start_index + height - 1) < #popup_table then
        start_index = start_index + 1
      end
    elseif dir == "up" then
      if cursor_pos > 1 then
        cursor_pos = cursor_pos - 1
      elseif start_index > 1 then
        start_index = start_index - 1
      end
    end
    render()
  end

  local function select_item()
    local index = start_index + cursor_pos - 1
    if popup_table[index] and callback then
      callback(popup_table[index])
    end
  end

  vim.keymap.set('n', '<ESC>', core.close_all_floats, { buffer = buffer })
  vim.keymap.set('n', '<Down>', function() move_cursor("down") end, { buffer = buffer })
  vim.keymap.set('n', '<Up>', function() move_cursor("up") end, { buffer = buffer })
  vim.keymap.set('n', '<CR>', select_item, { buffer = buffer })

  render()
  return win
end

return popup_menu

