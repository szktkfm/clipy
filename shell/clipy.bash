# https://github.com/szktkfm/clipy 
# Run daemon
clipy

__clipy_history() {
    local clipy_output
    clipy_output=$(clipy history --tmux)  
    READLINE_LINE="${READLINE_LINE}${clipy_output}"  
    READLINE_POINT=${#READLINE_LINE} 
}

# Bind the function to Ctrl+j
bind -x '"\C-j": clipy_history'
