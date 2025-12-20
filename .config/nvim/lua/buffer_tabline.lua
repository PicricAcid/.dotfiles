local api = vim.api

vim.cmd("highlight TabLine gui=underline guibg=none guifg=DarkCyan")
vim.cmd("highlight TabLineSel gui=underline guibg=DarkCyan guifg=black")
vim.cmd("highlight TabLineFill gui=underline guibg=none guifg=DarkCyan")

vim.o.tabline = "%!v:lua.BufferTabLine()"

function _G.BufferTabLine()
    local buffer_tabline = ""
    local sep = "|"
    
    local buffers_num = api.nvim_list_bufs()
    for _, buf_num in ipairs(buffers_num) do
	-- 通常のバッファのみ表示（ポップアップなどの特殊バッファを除外）
	local is_loaded = api.nvim_buf_is_loaded(buf_num)
	local is_listed = api.nvim_buf_get_option(buf_num, "buflisted")
	local buftype = api.nvim_buf_get_option(buf_num, "buftype")

	if is_loaded and is_listed and buftype == "" then
	    --- 選択しているバッファである場合、ハイライトを変える
	    local current_buf_num = api.nvim_get_current_buf()
	    
	    if buf_num == current_buf_num then
		buffer_tabline = buffer_tabline .. "%#TabLineSel#"
	    else
		buffer_tabline = buffer_tabline .. "%#TabLine#"
	    end

	    --- バッファ番号、バッファ名を表示
	    local buf_name = vim.fn.fnamemodify(api.nvim_buf_get_name(buf_num), ":t")
	    buffer_tabline = buffer_tabline .. buf_num .. ":" .. buf_name

	    --- バッファに変更が加えられている場合、"+"を表示
	    local modified = api.nvim_buf_get_option(buf_num, "modified")

	    if modified then
		buffer_tabline = buffer_tabline .. "+"
	    end

	    buffer_tabline = buffer_tabline .. sep
	end
    end

    buffer_tabline = buffer_tabline .. "%#TabLineFill#%T" .. "%=buffers"

    return buffer_tabline
end

