local api = vim.api
local data = require("comment_selector.data_loader")
local popup_menu = require("comment_selector.popup_menu")
local core = require("comment_selector.core")

local templates_win = 0
local categories_win = 0
local words_win = 0
local ok_win = 0
local comment_message = ""

local M = {}

local function insert_comment(win, message)
  api.nvim_set_current_win(win)
  local cursor_pos = api.nvim_win_get_cursor(0)
  local buf = api.nvim_get_current_buf()
  local commented = "-- " .. message
  api.nvim_buf_set_lines(buf, cursor_pos[1], cursor_pos[1], false, { commented })
end

function M.comment_menu()
  local origin_win = api.nvim_get_current_win()
  local buf = api.nvim_create_buf(false, true)
  api.nvim_buf_set_option(buf, 'bufhidden', 'wipe')
  api.nvim_buf_set_option(buf, 'modifiable', false)

  local message_buf = api.nvim_create_buf(false, true)
  local row, col = core.get_center_position(67, 16)

  local win = api.nvim_open_win(buf, true, {
    relative = 'editor',
    row = row,
    col = col,
    width = 67,
    height = 16,
    focusable = true,
    border = 'rounded',
    title = 'メッセージ入力',
    title_pos = 'left',
    noautocmd = true,
    zindex = 10,
  })

  api.nvim_win_set_option(win, 'winhighlight', 'FloatBorder:DarksoulsMenuBorder,NormalFloat:DarksoulsMenuText')

  local message_win = api.nvim_open_win(message_buf, true, {
    relative = 'win',
    row = 0,
    col = 0,
    width = 66,
    height = 1,
    focusable = true,
    border = 'none',
    noautocmd = true,
    zindex = 20,
    win = win,
  })

  local ok_buf = api.nvim_create_buf(false, true)
  api.nvim_buf_set_lines(ok_buf, 0, -1, true, { "OK<Enter>" })
  ok_win = api.nvim_open_win(ok_buf, true, {
    relative = 'win',
    row = 14,
    col = 45,
    width = 9,
    height = 1,
    focusable = true,
    border = 'rounded',
    noautocmd = true,
    zindex = 20,
    win = win,
    style = "minimal",
  })

  vim.keymap.set('n', '<ESC>', core.close_all_floats, { buffer = ok_buf })
  vim.keymap.set('n', '<CR>', function()
    insert_comment(origin_win, comment_message)
    core.close_all_floats()
  end, { buffer = ok_buf })

  templates_win = popup_menu(data.templates, {
    start_index = 1, relative = 'win', row = 2, col = 1, width = 20, height = 10,
    border = 'rounded', title = '定型文', zindex = 20, win = win,
  }, function(result)
    comment_message = result
    api.nvim_set_current_win(message_win)
    api.nvim_buf_set_lines(message_buf, 0, 1, false, { comment_message })
    api.nvim_set_current_win(categories_win)
  end)

  categories_win = popup_menu(data.categories, {
    start_index = 1, relative = 'win', row = 2, col = 23, width = 20, height = 10,
    border = 'rounded', title = 'カテゴリー', zindex = 20, win = win,
  }, function(result)
    local words = {}
    if result == "モノ" then
      words = data.words_object
    elseif result == "シチュエーション" then
      words = data.words_situation
    elseif result == "指示語" then
      words = data.words_directive
    elseif result == "属性" then
      words = data.words_type
    end

    if words_win ~= 0 and api.nvim_win_is_valid(words_win) then
      api.nvim_win_close(words_win, false)
    end

    words_win = popup_menu(words, {
      start_index = 1, relative = 'win', row = 2, col = 45, width = 20, height = 10,
      border = 'rounded', title = '単語', zindex = 20, win = win,
    }, function(word)
      if comment_message then
        comment_message = string.gsub(comment_message, "%*%*%*", word)
        api.nvim_set_current_win(message_win)
        api.nvim_buf_set_lines(message_buf, 0, 1, false, { comment_message })
        api.nvim_set_current_win(ok_win)
      else
        print("template is not selected")
      end
    end)
  end)

  api.nvim_set_current_win(templates_win)
end

return M

