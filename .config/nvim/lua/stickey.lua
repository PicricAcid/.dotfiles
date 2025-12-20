local M = {}
local api = vim.api

vim.cmd("highlight Stickey guifg=#000000 guibg=#ffff88")

function M.create_note(text)
    local buf = api.nvim_create_buf(false, true)

    local width = 0
    local lines = vim.split(text, "\n")
       
    local max_width = 0
    for _, line in ipairs(lines) do
	local width = vim.fn.strdisplaywidth(line)
	if width > max_width then
	max_width = width
	end
    end

    local size = math.max(max_width, #lines)

    while #lines < size do
	table.insert(lines, "")
    end

    for i, line in ipairs(lines) do  
	local padding = size - vim.fn.strdisplaywidth(line)
	lines[i] = line .. string.rep(" ", padding)
    end

    api.nvim_buf_set_lines(buf, 0, -1, false, vim.split(text, "\n"))

    local opts = {
	relative = 'editor',
	row = 3,
	col = 10,
	width = size + 2,
	height = size/2,
	style = 'minimal',
	border = 'rounded',
    }

    local win = api.nvim_open_win(buf, false, opts)
    api.nvim_win_set_option(win, 'winhighlight', 'FloatBorder:Stickey,NormalFloat:Stickey,Whitespace:Stickey')
    api.nvim_buf_set_option(buf, 'modifiable', true)
    
    api.nvim_create_autocmd('WinClosed', {
	callback = function(args)
	    if tonumber(args.match) == win then
		api.nvim_buf_delete(buf, { force = ture })
	    end
	end
    })
end

api.nvim_create_user_command('Memo', function(opts)
    M.create_note(opts.args)
end, { nargs = "+" })

return M
