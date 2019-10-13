# shellcheck disable=SC2164
cd "$(dirname $0)"
scp -4 debroot/*.deb lakka.kapsi.fi:public_html/debs/stretch
ssh -4 lakka.kapsi.fi bin/update_debs.sh
