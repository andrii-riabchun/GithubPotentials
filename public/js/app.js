/// <reference path="typings/jquery.d.ts" />
/// <reference path="typings/chart.d.ts" />

(function (){
       
function fillTableBody(element, dataArray) {
    var arr = [];
    for (var i in dataArray){
        var item = dataArray[i];
        arr.push("<tr><td>");
        arr.push(parseInt(i)+1);
        arr.push("</td><td><a href='https://github.com/");
        arr.push(item.fullname);
        arr.push(item.name);
        arr.push("'>");
        arr.push(item.fullname);
        arr.push(item.name);
        arr.push("</a>");
        arr.push("</td><td>");
        arr.push(item.value);
        arr.push("</td></tr>");          
    }
    element.html(arr.join(""));
}

function createErrorMessage(status){
    var result = "<h4>Oops! Something happened</h4><p>";
    switch (status) {
        case 400:
        case 404:
            result+="Seem like your request is incompatible with this server.";
            break;
        case 500:
            result+="Your request tried to break this well protected server.";
            break;
        case 418:
            result+="Github API limit rate exceeded. Come back in a hour!";
            break;
        case 0:
            result+="Server is gone. =(";
            break;
        default:
            result+="Unexpected error with code: "+status;
            break;
    }
    result+="</p>";
    return result;
}

function fillGraph(ctx, data){
    var labels = [];
    for (var i=data.days;i>0;i--) {
        labels.push(subDays(i));
    }
    var graphData = {
		labels : labels,
		datasets : [
			{   
                //#6f5499             
                label: "Stars",               
				fillColor : "rgba(111,84,153,0.5)",
				strokeColor : "rgba(111,84,153,0.8)",
				highlightFill: "rgba(111,84,153,0.7)",
				highlightStroke: "rgba(111,84,153,1)",
				data : data.starsdata
			},
            {
                //#996f54
                label: "Commits",
				fillColor : "rgba(153,111,84,0.5)",
				strokeColor : "rgba(153,111,84,0.8)",
				highlightFill: "rgba(153,111,84,0.7)",
				highlightStroke: "rgba(153,111,84,1)",
				data : data.commitsdata
			},
            {
                //#54996f
                label: "Contibutors",
				fillColor : "rgba(84,153,111,0.5)",
				strokeColor : "rgba(84,153,111,0.8)",
				highlightFill: "rgba(84,153,111,0.7)",
				highlightStroke: "rgba(84,153,111,1)",
				data : data.contribsdata
			},
		]
	};
    

    return new Chart(ctx).Bar(graphData, {
			    responsive : true
		    });
}

function fillInfoTable(data) {    
    var owner = "<a href='https://github.com/"+
        data.owner+"'>"+
        data.owner+"</a>";
    var repo = "<a href='https://github.com/"+
        data.owner+"/"+data.repo+"'>"+
        data.repo+"</a>";
    
    var arr = [];
    arr.push("<tr><td>Owner:</td><td>");
    arr.push(owner);
    arr.push("</td></tr><tr><td>Repo:</td><td>");
    arr.push(repo);
    arr.push("</td></tr><tr><td>Stars:</td><td>");
    arr.push(data.stars);
    arr.push("</td></tr><tr><td>Commits:</td><td>");
    arr.push(data.commits);
    arr.push("</td></tr><tr><td>Contribs:</td><td>");
    arr.push(data.contribs);
    arr.push("</td></tr>");
    
    // <tr>
    //     <td>Owner:</td>
    //     <td id="repoinfo-owner"></td>    
    // </tr>
    // <tr>
    //     <td>Repo:</td>
    //     <td id="repoinfo-repo"></td>    
    // </tr>
    // <tr>
    //     <td>Stars:</td>
    //     <td id="repoinfo-stars"></td>    
    // </tr>
    // <tr>
    //     <td>Commits:</td>
    //     <td id="repoinfo-commits"></td>    
    // </tr>
    // <tr>
    //     <td>Contributors:</td>
    //     <td id="repoinfo-contribs"></td>    
    // </tr>

    $("#repoinfo-info").html(arr.join(""));
}

function subDays(value) {
    var d = new Date();
    d.setDate(d.getDate()-value);
    var day =   d.getDate().toString();
    var month = d.getMonth().toString();
    if (day.length===1)		day=0+day;
    if (month.length===1)   month=0+month;
    return day+"."+month;   
}

$("#repos-button").click(function(event){
    table = $("#repos-table");
    button = $("#repos-button");
    spinner = $("#repos-spinner");
    callout = $("#repos-callout");
       
    button.prop("disabled", true);       
    spinner.show();      
    callout.hide();
    
    
    var criteria = $("#repos-criteria-select").val();
    var timespan = $("#repos-timespan").val();
    var targetUrl = "/repos/"+criteria+"/"+timespan;
            
    $.ajax({
        url: targetUrl,
        dataType: 'json',                         
        success: function(data){                
            fillTableBody($("#repos-table-body"), data.items);
            spinner.hide();
            button.prop("disabled", false);
            table.show();
        },
        error: function (response) {           
            var message = createErrorMessage(response.status);
            callout.html(message);
            callout.show();
            spinner.hide();
            button.prop("disabled", false);     
        }           
        });
});
    
$("#orgs-button").click(function(event){
    table = $("#orgs-table");
    button = $("#orgs-button");
    spinner = $("#orgs-spinner");
    callout = $("#orgs-callout");
         
    button.prop("disabled", true);       
    spinner.show();      
    callout.hide();
    
    var criteria = $("#orgs-criteria-select").val();
    var timespan = $("#orgs-timespan").val();
    var targetUrl = "/orgs/"+criteria+"/"+timespan;
            
    $.ajax({
        url: targetUrl,
        dataType: 'json',                         
        success: function(data){               
            fillTableBody($("#orgs-table-body"), data.items);
            table.show();
            spinner.hide();
            button.prop("disabled", false); 
        },
        error: function (response) {
            var message = createErrorMessage(response.status);
            callout.html(message);
            callout.show();
            spinner.hide();
            button.prop("disabled", false);
        }           
        });
});   

$("#repoinfo-button").click(function(event){
    button = $("#repoinfo-button");
    spinner = $("#repoinfo-spinner");
    callout = $("#repoinfo-callout");
    canvas = $("#repoinfo-canvas");
    container = $("#repoinfo-container");
    ctx = canvas.get(0).getContext("2d");
         
    button.prop("disabled", true);       
    spinner.show();      
    callout.hide();
    container.hide();
    
    
    var name = $("#repoinfo-fullname").val();
    var timespan = $("#repoinfo-timespan").val();
    var targetUrl = name+"/"+timespan;
            
    $.ajax({
        url: targetUrl,
        dataType: 'json',                         
        success: function(data){
            container.show();
            if (window.graph !== undefined) window.graph.destroy();           
            window.graph = fillGraph(ctx,data);
            fillInfoTable(data);
            spinner.hide();
            button.prop("disabled", false);
        },
        error: function (response) {
            container.show();
            var message = createErrorMessage(response.status);
            callout.html(message);
            callout.show();
            spinner.hide();
            button.prop("disabled", false);
        }           
        });
}); 
  
})();