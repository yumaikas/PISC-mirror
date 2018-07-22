" syntax/pisc.vim

" PISC only has a few 'keywords', so we'll match them here.
syntax keyword piscKeyWords if when while break continue :PRE :DOC 


" : name 
" NOTE: there are tabs and spaces in those gaps
syntax match piscFuncName /\v:[ \t]+[^ \t]+[     ]+/

" Match ;
syntax match piscKeyWords "\v\;"

" Match basic PISC numbers
" There are both tabs and spaces in both of the (||$) sections
syntax match piscNumber /\v(^| |\t)\@=-?\d+(  | |$)\@=/
syntax match piscNumber /\v(^| |\t)\@=-?\d+\.\d*( | |$)\@=/
syntax match piscBool "\v<t>"
syntax match piscBool "\v<f>" 

syntax match piscPrefixWord /\v(^| |\t)\@=[-_:!@$%^&<>+?\.]+[-a-zA-Z0-9]+>/

syntax match piscBraces "\v\$\{"
syntax match piscBraces "\v\{"
syntax match piscBraces "\v\}"
syntax match piscBraces "\v\["
syntax match piscBraces "\v\]"

" Strings
syntax region piscString start=/"/ skip=/\v\\./ end=/"/

" Comments
syntax region piscLineComment start=/\v#/ end=/\v$/ oneline
syntax region piscMultilineComment start=/\v\/\*/ end=/\v\*\//
syntax region piscStackComment start=/\v\(/ end=/\v\)/


" Set highlights
highlight default link piscNumber Number
highlight default link piscKeyWords Statement
highlight default link piscBool Boolean
highlight default link piscString String
highlight default link piscBraces Statement
highlight default link piscPrefixWord Identifier
highlight default link piscFuncName Statement

highlight default link piscLineComment Comment
highlight default link piscMultilineComment Comment
highlight default link piscStackComment Comment

