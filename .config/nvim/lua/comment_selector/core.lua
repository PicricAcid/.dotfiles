local api = vim.api

vim.cmd("highlight DarksoulsFloatBorder guifg=#007acc")
vim.cmd("highlight DarksoulsFloatTitle guifg=#005f87")
vim.cmd("highlight DarksoulsFloatText guifg=#4a4a4a guibg=#ffffff")
vim.cmd("highlight DarksoulsFloatTextSelected guifg=#ffffff guibg=#007acc")
vim.cmd("highlight DarksoulsMenuBorder guifg=#007acc")
vim.cmd("highlight DarksoulsMenuText guifg=#4a4a4a guibg=#ffffff")

local M = {}

function M.get_center_position(width, height)
  local ui = api.nvim_list_uis()[1]
  local row = math.floor((ui.height - height) / 2)
  local col = math.floor((ui.width - width) / 2)
  return row, col
end

function M.close_all_floats()
  for _, win in ipairs(api.nvim_list_wins()) do
    if api.nvim_win_get_config(win).relative ~= "" then
      api.nvim_win_close(win, true)
    end
  end
end

return M

