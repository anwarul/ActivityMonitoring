const $rowTemplate = $($("#rowTemplate").html())
const $mainTable = $("#main-table")

function deleteDevice(deviceAddr){
	$.ajax({
		url: "/api/deletedevice?physaddr="+deviceAddr,
		dataType: "json",
		success: (data) => {
			updateActivityList();
			alert(JSON.stringify(data));
		}
	});
}

function renderActivity(result){
	for(var mac in result){
		const value = result[mac]
		$row = $rowTemplate.clone();
		$row.find(".mac").text(value.PhysicalAddress);
		$row.find(".name").text(value.Name);
		$row.find(".last-active").text(value.LastResponse);
		$row.find(".is-enabled").text(value.IsActive);
		$row.find(".delete-button").click(()=>{
			deleteDevice(value.PhysicalAddress);
		});
		$mainTable.append($row)
	}
}

function updateActivityList() {
	$(".tableRow").remove();
	$.ajax({
		url: "/api/devicestatus",
		dataType: "json",
		success: renderActivity
	});
}

function getTickTime(){
	$.ajax({
		url: "/ticktime",
		dataType: "text",
		success: (data) => {
			$("#tick-time").text(data);
		}
	});
}

$(document).ready(function(){
	getTickTime();
	updateActivityList();
	setInterval(updateActivityList, 1000*10); //10 seconds
})