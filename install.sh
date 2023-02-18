HERE=$PWD

COPY=1 # false

while getopts ":p" opt; do
    case $opt in
        p) # preserve files in repo directory
            COPY=0 # true
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

move_file() {
    if [ $1 = 0 ]; then
        cp $2 $3
    else
        mv $2 $3
    fi
}

move_dir() {
    if [ $1 = 0 ]; then
        cp -r $2 $3
    else
        mv -r $2 $3
    fi
}

USRBIN="$HOME/bin"

if [ ! -d "$USRBIN" ]; then
    mkdir "$USRBIN"
fi

KNOTBIN="$HERE/bin/knot"

echo "moving binary from ${KNOTBIN} to ${USRBIN}"

move_file "$COPY" "$KNOTBIN" "$USRBIN"

KNOTCONFIG="$HOME/.config/knot"

if [ ! -d $KNOTCONFIG ]; then
    echo "moving config info to ${KNOTCONFIG}"

    mkdir $KNOTCONFIG

    PROJECTS="$HERE/projects.json"

    TEMPLATES="$HERE/templates"

    if [ -f $PROJECTS ]; then
        move_file $COPY $PROJECTS $KNOTCONFIG
    else
        touch "$KNOTCONFIG/projects.json"
    fi

    if [ -d $TEMPLATES ]; then
        move_dir $COPY $TEMPLATES $KNOTCONFIG
    else
        mkdir "$KNOTCONFIG/templates"
    fi

else
    echo "cannot use directory ${KNOTCONFIG} as config because it already exists"
fi

