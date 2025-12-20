local api = vim.api
local eval = require("comment_selector.evaluation")

local eval_win = nil

local function show_evaluation_popup()
  local line = api.nvim_win_get_cursor(0)[1]
  local result = eval.get_evaluation(line)

  if result then
    local text = string.format("評価: %d", result.score)

    if eval_win and api.nvim_win_is_valid(eval_win) then
      api.nvim_win_close(eval_win, true)
    end

    local buf = api.nvim_create_buf(false, true)
    api.nvim_buf_set_lines(buf, 0, -1, false, { text })

    eval_win = api.nvim_open_win(buf, false, {
      relative = "cursor",
      row = 1,
      col = 0,
      width = #text + 2,
      height = 1,
      style = "minimal",
      border = "rounded",
    })

    api.nvim_win_set_option(eval_win, 'winhighlight', 'NormalFloat:DarksoulsFloatText')
  else
    if eval_win and api.nvim_win_is_valid(eval_win) then
      api.nvim_win_close(eval_win, true)
      eval_win = nil
    end
  end
end

vim.api.nvim_create_autocmd("CursorMoved", {
  callback = show_evaluation_popup,
})

