# clipy
function __clipy_history() {
    BUFFER="${BUFFER}"$(/tmp/clipy history)
    CURSOR=${#BUFFER}  
    zle redisplay
}
zle -N __clipy_history
bindkey "^j" __clipy_history
#       ^^^ 
#       ctrl + j 
