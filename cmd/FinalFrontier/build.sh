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

# For now, force re-packing of bindata
rm bindata.go

older_than bindata.go strings/*.pisc && {
  echo "Packing scripts"
  go-bindata -pkg main -o bindata.go scripts/...
}

case "$1" in
  build)
  older_than FinalFrontier *.go && {
    go build -o FinalFrontier
  }
  ;;
  run)
  older_than FinalFrontier main.go finalFrontier.go && {
    rm FinalFrontier
    go build -o FinalFronter
  }
  ./FinalFrontier
  ;;
esac
  




