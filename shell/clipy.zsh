# https://github.com/szktkfm/clipy 
# Run daemon
clipy

__clipy_history() {
    BUFFER="${BUFFER}"$(clipy history --tmux)
    CURSOR=${#BUFFER}  
    zle redisplay
}
zle -N __clipy_history

# Bind the function to Ctrl+j
bindkey "^j" __clipy_history
