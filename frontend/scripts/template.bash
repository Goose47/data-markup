#!/bin/bash

set -e

SCRIPT_DIR=$(dirname "$(realpath "$0")")
COMPONENT_NAME=$1
NAMESPACE_NAME=$(echo "$COMPONENT_NAME" | awk '{
    for (i=1; i<=length($0); i++) {
        c = substr($0, i, 1);
        if (c ~ /[A-Z]/ && i > 1) {
            printf "-%s", tolower(c);
        } else {
            printf "%s", tolower(c);
        }
    }
    printf "\n";
}')
COMPONENT_PATH="$SCRIPT_DIR/../src/components/$COMPONENT_NAME"
SCSS_FILE="$SCRIPT_DIR/../src/components/$COMPONENT_NAME/$COMPONENT_NAME.scss"
TSX_FILE="$SCRIPT_DIR/../src/components/$COMPONENT_NAME/$COMPONENT_NAME.tsx"

mkdir $COMPONENT_PATH
touch $SCSS_FILE
touch $TSX_FILE

echo "@use '../../variables.scss';" >> $SCSS_FILE
echo "\$block: '.#{variables.\$ns}$NAMESPACE_NAME';" >> $SCSS_FILE
echo '' >> $SCSS_FILE
echo '#{$block} {' >> $SCSS_FILE
echo '' >> $SCSS_FILE
echo '}' >> $SCSS_FILE

echo "import { block } from \"../../utils/block\";" >> $TSX_FILE
echo "import \"./$COMPONENT_NAME.scss\";" >> $TSX_FILE
echo '' >> $TSX_FILE
echo "const b = block(\"$NAMESPACE_NAME\");" >> $TSX_FILE
echo '' >> $TSX_FILE
echo "export const $COMPONENT_NAME = () => {" >> $TSX_FILE
echo "return (<div className={b()}></div>)" >> $TSX_FILE
echo "};" >> $TSX_FILE

echo "Done. $COMPONENT_NAME created in $COMPONENT_PATH"
