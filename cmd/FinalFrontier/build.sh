#! /bin/bash
function older_than() {
  local target="$1"
  shift
  if [ ! -e "$target" ]
  then
    echo "updating $target" >&2
    return 0  # success
  fi
  local f
  for f in "$@"
  do
    if [ "$f" -nt "$target" ]
    then
      echo "updating $target" >&2
      return 0  # success
    fi
  done
  return 1  # failure
}

older_than bindata.go strings/*.pisc && {
  echo "Packing scripts"
  go-bindata -pkg main -o bindata.go scripts/...
}

case "$1" in
  build)
  older_than game.exe *.go && {
    go build -o game.exe
  }
  ;;
  run)
  older_than game.exe main.go finalFrontier.go && {
    rm game.exe
    go build -o game.exe
  }
  ./game.exe
  ;;
esac
  




