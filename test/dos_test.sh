#Make sure prometheus server url is specified
test -n "$PROMETHEUS_URL" || eval 'echo "Need to specify a PROMETHEUS_URL enviroment variable where to collect metrics" ; exit'

#Need to clone GoldenEye into parent folder of kdd
#pushd ../../GoldenEye

current_unix_timestamp=$(date '+%s')
echo "Current timestamp: $current_unix_timestamp"

times=300
end_time=$current_unix_timestamp
start_time=$((end_time-times))
step="1"
echo "Start time       : $start_time"
echo "End time         : $end_time"

command="curl '$PROMETHEUS_URL/api/v1/query_range?query=cilium_drop_bytes_total&step=$step&start=$start_time&end=$end_time' | jq"
echo "Command to be execute:"
echo "$command"

metric=$(eval "$command")
#echo "$metric"

need_to_watch=$(echo "$metric" | awk "/\"metric\":/{want=0} /Policy denied by denylist/{want=1} want")
#This line for testing command to know it works
#need_to_watch=$(echo "$metric" | awk "/\"metric\":/{want=0} /\"reason\": \"Stale or unroutable IP\"/{want=1} want") 
[[ -z "$need_to_watch" ]] && eval 'echo "No metric for cilium drop via denied by policy" ; exit'

#echo "$need_to_watch"
values=$(echo "$need_to_watch" | awk "/\[/{point=point+1} /\]/{point=point-1} /value/{next} point>0")

#echo "$values"

output=$(echo "$values" | \
awk -v sum="$sum" -v periods="$periods" \
-F"\"" "/\"[0-9]*\"/{this=\$2;} prev > 0 && prev < this && //{sum=sum+this-prev; periods=periods+1} //{prev=this} END {print sum; print period>
)
sum=$(echo "$output" | head -1)
periods=$(echo "$output" | tail -1)

echo "sum $sum periods $periods"

[[ periods -eq 0 ]] && eval 'echo "No surge of drop traffic" ; exit'

average=$((sum*8/periods/1000))
echo "$total_in_periods"
echo "The average: $average Kbps, in $times s from $start_time to $end_time"

#popd
