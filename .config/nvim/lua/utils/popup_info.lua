local M = {}
local api = vim.api

vim.cmd("highlight PopupInfoFloatBorder guifg=#006db3")
vim.cmd("highlight PopupInfoFloatTitle guifg=#6ab7ff")
vim.cmd("highlight PopupInfoText guifg=#abb2bf")

local popup_state = {
  buffer = nil,
  window = nil,
  timer = nil
}

local function close_popup()
  if popup_state.timer then
    vim.fn.timer_stop(popup_state.timer)
    popup_state.timer = nil
  end

  if popup_state.window and api.nvim_win_is_valid(popup_state.window) then
    api.nvim_win_close(popup_state.window, true)
    popup_state.window = nil
  end

  if popup_state.buffer and api.nvim_buf_is_valid(popup_state.buffer) then
    api.nvim_buf_delete(popup_state.buffer, { force = true })
    popup_state.buffer = nil
  end
end

---@param text string|table Text to display (string or table of lines)
---@param opt table|nil Options for popup window
---@return nil
function M.popup_info(text, opt)
  close_popup()

  local lines = type(text) == "string" and { text } or text

  if opt == nil then
    opt = {
      relative = 'cursor',
      row = 0,
      col = 0,
      width = 20,
      height = 1,
      border = 'rounded',
      title = 'Info',
      zindex = 200,
      timeout = 3000,
    }
  end

  local width = opt.width or 20
  local height = opt.height or math.min(#lines, 10)

  if not opt.width then
    for _, line in ipairs(lines) do
      width = math.max(width, #line + 2)
    end
  end

  local title = opt.title and { { opt.title, 'PopupInfoFloatTitle' } } or nil

  local ui = api.nvim_list_uis()[1]
  if not ui then
    return
  end

  local row = opt.row or 0
  local col = opt.col or ui.width - width - 1
  local relative = opt.relative or 'editor'

  if relative == 'cursor' then
    col = opt.col or 0
  end

  popup_state.buffer = api.nvim_create_buf(false, true)

  api.nvim_buf_set_option(popup_state.buffer, 'buftype', 'nofile')
  api.nvim_buf_set_option(popup_state.buffer, 'bufhidden', 'wipe')
  api.nvim_buf_set_option(popup_state.buffer, 'swapfile', false)
  api.nvim_buf_set_option(popup_state.buffer, 'buflisted', false)

  api.nvim_buf_set_lines(popup_state.buffer, 0, -1, false, lines)

  popup_state.window = api.nvim_open_win(popup_state.buffer, false, {
    relative = relative,
    row = row,
    col = col,
    width = width,
    height = height,
    focusable = false,
    border = opt.border or 'rounded',
    title = title,
    title_pos = 'left',
    noautocmd = true,
    zindex = opt.zindex or 200,
  })

  api.nvim_win_set_option(popup_state.window, 'number', false)
  api.nvim_win_set_option(popup_state.window, 'relativenumber', false)
  api.nvim_win_set_option(popup_state.window, 'wrap', false)
  api.nvim_win_set_option(popup_state.window, 'cursorline', false)
  api.nvim_win_set_option(popup_state.window, 'winhighlight', 'FloatBorder:PopupInfoFloatBorder,NormalFloat:PopupInfoText')

  if opt.timeout and opt.timeout > 0 then
    popup_state.timer = vim.fn.timer_start(opt.timeout, function()
      close_popup()
    end)
  end
end

function M.close()
  close_popup()
end

return M
