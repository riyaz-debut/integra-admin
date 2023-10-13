# #!/bin/sh

ADD_ORG_ENVELOPE_JSON="$1"
ADD_ORG_ENVELOPE_PB="$2"







configtxlator proto_encode --input ${ADD_ORG_ENVELOPE_JSON} --type common.Envelope > ${ADD_ORG_ENVELOPE_PB}















# # for std out
# jq -s '.[0] * {"channel_group":{"groups":{"Application":{"groups": {"'${ORG_NAME}'":.[1]}}}}}' ${SOURCE} ${ORG_SOURCE}


# # for out to file
# jq -s '.[0] * {"channel_group":{"groups":{"Application":{"groups": {"'${ORG_NAME}'":.[1]}}}}}' ${SOURCE} ${ORG_SOURCE} > ${FINAL}