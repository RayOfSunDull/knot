KNOTDIR=~/projects/code/programs/knot

go build -o $KNOTDIR/bin/knot $KNOTDIR/src

while getopts ":r" opt; do
    case $opt in
        r) # run
            $KNOTDIR/bin/knot
            ;;
        \?)
			echo "invalid option: -$OPTARG"
			exit 1
			;;
		*)
			echo "option -$OPTARG requires an argument"
			exit 1
			;;
    esac
done
