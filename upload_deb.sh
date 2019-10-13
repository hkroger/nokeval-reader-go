# shellcheck disable=SC2164
cd "$(dirname $0)"
scp debroot/*.deb lakka.kapsi.fi:public_html/debs/stretch
ssh lakka.kapsi.fi bin/update_debs.sh
