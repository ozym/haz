/***
 * Date functions
 * ***/
(function (Date, undefined) {
    var origParse = Date.parse, numericKeys = [ 1, 4, 5, 6, 7, 10, 11 ];
    var nzdtStartCache = {}; //cache nzdt start/end date by year to reduce loops
    var nzdtEndCache = {};

    var padNum = function (number, length) {
        var str = '' + number;
        while (str.length < length) {
            str = '0' + str;
        }
        return str;
    };

    Date.parseUTC = function (date) {
        var timestamp, struct, minutesOffset = 0;
        /**
        * S5 §15.9.4.2 states that the string should attempt to be parsed as a Date Time String Format string
        * before falling back to any implementation-specific date parsing, so that’s what we do, even if native
        * implementations could be faster
        * 1 YYYY 2 MM 3 DD 4 HH 5 mm 6 ss 7 msec 8 Z 9 ± 10 tzHH 11 tzmm
        * "2011-12-08 09:48:19.000000"
        *
        * 2012-03-27 00:09:44.074000
        * /^(\d{4})-(\d{2})-(\d{2}) (\d{2}):(\d{2}):(\d{2}).(\d{1,10})$/        *
        */
        //if ((struct = /^(\d{4}|[+\-]\d{6})(?:-(\d{2})(?:-(\d{2}))?)?(?:T(\d{2}):(\d{2})(?::(\d{2})(?:\.(\d{3}|[+\-]))?)?(?:(Z)|([+\-])(\d{2})(?::(\d{2}))?)?)?$/.exec(date))) {
        if ((struct = /^(\d{4})-(\d{1,2})-(\d{1,2})T(\d{1,2}):(\d{1,2}):(\d{1,2}).(\d{1,10})$/.exec(date))) {
            // avoid NaN timestamps caused by "undefined" values being passed to Date.UTC
            for (var i = 0, k; (k = numericKeys[i]); ++i) {
                struct[k] = +struct[k] || 0;
            }
            // allow undefined days and months
            struct[2] = (+struct[2] || 1) - 1;
            struct[3] = +struct[3] || 1;

            timestamp = Date.UTC(struct[1], struct[2], struct[3], struct[4], struct[5], struct[6]);
        }
        else {
            timestamp = origParse ? origParse(date) : NaN;
        }
        return new Date(timestamp);
    };

    Date.toNZTimeString = function (date) {
        var tz  ;
        var offset = Date.getNZTimeZone(date);
        if(offset == 13){
            tz = "NZDT";
        }
        if(offset == 12){
            tz = "NZST";
        }
        if(tz)	{
            tz = " (" + tz + ")"
        }else{
            tz = '';
        }
        //console.log("timezone " + tz);
        var date1 = new Date(date.getTime() + offset*60*60*1000);//add offset
        return date1.getUTCFullYear() + "-"	+ (padNum((date1.getUTCMonth() +1),2)) + "-" + (padNum(date1.getUTCDate(),2))
        + " " + (padNum(date1.getUTCHours(),2)) + ":" + (padNum(date1.getUTCMinutes(),2)) + ":" + (padNum(date1.getUTCSeconds(),2))  + tz ;
    };
    Date.toUTCString = function (date) {
        return date.getUTCFullYear() + "-"	+ (padNum((date.getUTCMonth() +1),2)) + "-" + (padNum(date.getUTCDate(),2))
        + " " + (padNum(date.getUTCHours(),2)) + ":" + (padNum(date.getUTCMinutes(),2)) + ":" + (padNum(date.getUTCSeconds(),2)) + " (GMT)";
    };

    /*
	 * NZ DST begins at 02:00 NZST on the last Sunday in September
	 */
    Date.getNZDTStart = function (date) {
        var year = date.getUTCFullYear();
        //get from cache first
        var nzdtStart = nzdtStartCache["" + year];
        //console.log("1>>nzdtStart " + nzdtStart);
        if(nzdtStart){
            return nzdtStart;
        }
        //console.log("2>>nzdtStart "  );
        nzdtStart = new Date(Date.UTC(year, 8, 29, 14,0,0));//end of sept nz time
        //get last sunday of sept
        while(nzdtStart.getUTCMonth() > 7){
            var nzdtStart1 = new Date(nzdtStart.getTime() + 12*60*60*1000);
            if(nzdtStart1.getUTCDay() == 0){//sunday
                break;
            }
            nzdtStart.setTime(nzdtStart.getTime() - 1*24*60*60*1000);//go to previous day
        }
        //cache it
        nzdtStartCache["" + year] = nzdtStart;
        return nzdtStart;
    };

    /*
	 * NZ DST ends at 03:00 NZDT
	 * (or 02:00 NZST as defined in the Time Act 1974) on the first Sunday in April.
	 */
    Date.getNZDTEnd = function (date) {
        var year = date.getUTCFullYear();
        //get from cache first
        var nzdtEnd = nzdtEndCache["" + year];
        //console.log("1>>nzdtEnd " + nzdtEnd);
        if(nzdtEnd){
            return nzdtEnd;
        }
        //console.log("2>>nzdtEnd " + nzdtEnd);
        nzdtEnd = new Date(Date.UTC(year, 2, 31,14,0,0)); //start of april nz time
        //get first sunday of april
        while(nzdtEnd.getUTCMonth() < 3                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                  ){
            var nzdtEnd1 = new Date(nzdtEnd.getTime() + 12*60*60*1000);
            if(nzdtEnd1.getUTCDay() == 0){//sunday
                break;
            }
            nzdtEnd.setTime(nzdtEnd.getTime() + 1*24*60*60*1000);//go to next day
        }
        nzdtEndCache["" + year] = nzdtEnd;
        return nzdtEnd;
    };
    //
    Date.getNZTimeZone = function(date){
        if(!date) {//use the current date
            date = new Date();
        }
        var nzdtStart = Date.getNZDTStart(date);
        var nzdtEnd =  Date.getNZDTEnd(date);

        if(date.getTime() > nzdtStart.getTime() || date.getTime() < nzdtEnd.getTime() ){
            return 13;// "NZDT";
        }
        else{
            return 12;//"NZST";
        }
    }
}(Date));

var qsInteractive = {
    MAX_NUMBER_QUAKES:20000,
    SERVICE_PATH : '',
    mapCentreNZ: [-41.27, 173.28],
    selectedRegion:null,
    regionCenters:{
        "newzealand": [-38.2057205720572, 175.788778877888],
        "aucklandnorthland": [-36.0897953716283, 174.130454598145],
        "tongagrirobayofplenty": [-37.9474981652609, 176.539583931146],
        "gisborne": [-38.0284634174816, 178.427911861592],
        "hawkesbay": [-39.6263376999659, 177.190468417066],
        "taranaki": [-39.1273951052499, 173.89201240672],
        "wellington": [-41.2177571660282, 175.293934996233],
        "nelsonwestcoast": [-42.192553054577, 170.933790737804],
        "canterbury": [-43.6338286051729, 172.715320498205],
        "fiordland": [-45.7607533091421, 166.30255175867],
        "otagosouthland": [-46.522978901575, 168.974427994776]
    },
    regionPolies:{
        "aucklandnorthland": [[-38.138, 173.251], [-38.045, 175.583 ], [-36.379, 176.474 ], [-34.026, 174.285 ], [-34.135, 171.857 ], [-38.138, 173.251 ]],
        "tongagrirobayofplenty": [[-39.526, 175.028 ], [-39.809, 175.722 ], [-38.688, 176.931 ], [-36.770, 178.346 ], [-36.379, 176.474 ], [-38.045, 175.583 ], [-39.526, 175.028 ]],
        "gisborne": [[-38.688, 176.931 ], [ -39.274, 178.561], [-37.361, 179.898 ], [-36.770, 178.346 ], [-38.688, 176.931 ]],
        "hawkesbay": [[-38.688, 176.931 ], [-39.809, 175.722 ], [-40.638, 177.560 ], [-39.274, 178.561 ], [-38.688, 176.931 ]],
        "taranaki": [[-39.632, 172.004 ], [-40.456, 174.156 ], [-39.526, 175.028 ], [-38.045, 175.583 ], [-38.138, 173.251 ], [-39.632, 172.004 ]],
        "wellington": [[-41.767, 172.951 ], [ -42.908, 175.748], [-40.638, 177.560 ], [-39.526, 175.028 ], [-40.462, 174.109 ], [-41.767, 172.951 ]],
        "nelsonwestcoast": [[-43.711,167.399 ], [-44.668, 169.168 ], [-44.341, 169.564 ], [-42.832, 172.001 ], [-41.767, 172.951 ], [-40.462, 174.109 ],
        [-39.632, 172.004 ], [-41.892, 170.180 ], [-43.711, 167.399 ]],
        "canterbury": [[-41.767,172.951 ], [-42.832, 172.001 ], [-44.341, 169.564 ], [-45.412, 172.312 ], [-42.908, 175.748 ], [-41.767, 172.951 ]],
        "fiordland": [[-46.083,164.218 ], [-47.212, 163.787 ], [-47.827, 165.247 ], [-44.668, 169.168 ], [-43.711, 167.399 ], [-46.083, 164.218 ]],
        "otagosouthland": [[-47.827, 165.247 ], [-48.410, 169.148 ], [-45.412, 172.312 ], [-44.341, 169.564 ], [-44.668, 169.168 ], [-47.827, 165.247 ]]
    },
    outputOptions : ["GeoJSON", "GML", "KML", "CSV"], //, "Image"
    imageSizes : ["512", "1024", "2048"],
    maxNumQuakes : ["100", "500", "1000", "1500", "2000"],
    locationByMap:true, //if true, the map extent will be populated in the location field
    resultOption:'query', // or map
    resultFormat:null, //query result format
    currentURL:null,
    resultDataUrl:null,
    resultMaxQuakes:null,
    resultCountUrl:null,
    queryString : "",  //query string for search result, excludes date fields
    queryStringDate : "",//query string for search result, date fields
    imageSize:512,
    imageShowPlaces:0,

    initDatePicker:	function (){
        //init date picker
        jQuery('.date').datepicker({
            format: 'yyyy-mm-dd'
        });
    },

    initEvents:function (){
        jQuery("input[name='dateRadios']").on('change', function () {
            if(jQuery("input[name='dateRadios']:checked") ){
                //console.log("dateRadios change " + jQuery("input[name='dateRadios']:checked") );
                var selectedVal = jQuery("input[name='dateRadios']:checked").val();
                qsInteractive.fillDateFields(selectedVal);
            }
            qsInteractive.resetMap();
        });

        jQuery("input[name='locRadios']").on('change', function () {
            var selectedVal = jQuery("input[name='locRadios']:checked").val();
            //console.log("selectedVal " + selectedVal);
            qsInteractive.resetMap();
            qsInteractive.setLocationOption(selectedVal,true);
        });

        jQuery("input[name='resultRadios']").on('change', function () {
            var selectedVal = jQuery("input[name='resultRadios']:checked").val();
            //console.log("resultRadios selectedVal " + selectedVal);
            qsInteractive.setResultOption(selectedVal);
        });
        //query output format chnge
        jQuery('#outputFormat').on('change',function(){
            var mapformat =  this.value.toLowerCase();
            if(mapformat == 'image'){
                jQuery("#image-options").removeClass("hide").addClass("inline");
            }else{
                jQuery("#image-options").removeClass("inline").addClass("hide");
            }
            qsInteractive.hideQueryResult();
        });

        jQuery('#btnDateClear').on('click', function (e) {
            qsInteractive.clearDateFields();
        });

        jQuery('#btnLocationClear').on('click', function (e) {
            qsInteractive.clearLocationFields(true);
        });
        //btnMagDepthQualityClear
        jQuery('#btnMagDepthQualityClear').on('click', function (e) {
            qsInteractive.clearMagDepthFields();
            qsInteractive.clearQualityFields();
        });

        jQuery('#btnShowQuery').on('click', function (e) {
            qsInteractive.getQuerySearchResult();
        });

        jQuery('#btnShowMap').on('click', function (e) {
            var val = jQuery(this).val();
            if(val =='Show'){
                qsInteractive.getMapSearchResult();
            }else{//clear
                qsInteractive.resetMap();
                if(qsInteractive.locationByMap){
                    //force map update coordinates fields
                    quakesMapApp.updateMapBound();
                }
            }
        });

        jQuery('#btnReset').on('click', function (e) {
            qsInteractive.resetFields();
        });

        jQuery("#mapTownCheck").on('change', function () {
            qsInteractive.mapShowPlaces = this.checked;
            qsInteractive.hideQueryResult();
        });

        jQuery("#imgSizeSelect").on('change', function () {
            qsInteractive.imageSize = jQuery("#imgSizeSelect").val();
            //qsInteractive.imageSize = this.val();
            qsInteractive.hideQueryResult();
        });

        //shown
        jQuery('#divLatLon').on('shown', function () {
            jQuery("#btnLatLon i").removeClass("icon-chevron-down").addClass("icon-chevron-up");
            jQuery("#btnLatLon").attr("title","Hide coordinate fields");
        });
        jQuery('#divLatLon').on('hidden', function () {
            qsInteractive.clearQualityFields();
            jQuery("#btnLatLon i").removeClass("icon-chevron-up").addClass("icon-chevron-down");
            jQuery("#btnLatLon").attr("title","Show coordinate fields");
        });

        //validate numeric fileds
        jQuery("input[name='numField']").on('keypress', function (e) {
            if (e.ctrlKey != true) {
                return qsInteractive.checkNum(e, [45, 46]);//46 = dot
            }
            return true;
        });
        //validate paste value
        jQuery("input[name='numField']").on('paste', function (e){
            var pasted = e.originalEvent.clipboardData.getData('Text');
            if(!jQuery.isNumeric( pasted )){
                setTimeout(function () {
                    e.target.value = '';
                }, 100);
            }
        });

        //validate date fileds
        jQuery(".date input[type='text']").on('keypress', function (e) {
            return qsInteractive.checkNum(e, [45]);//45 = -
        });

        //coordinat fields
        jQuery("input[name='numField']").on('click', function (e) {
            var currentId = jQuery(this).attr('id');
            if(currentId == 'minLat'||currentId == 'minLon'||currentId == 'maxLat'||currentId == 'maxLon'){
                jQuery("#locCustom")[0].checked = true; //force check the custom radio button
            }
        });

    },

    hideQueryResult:function(){
        jQuery("#result").css("display", "none");
        jQuery("#btnReset").css("display", "none");
    },

    fillMapBound:function (bbox){
        //console.log(" fillMapBound " + this.locationByMap);
        if(this.locationByMap){
            if(bbox){
                this.fillLocationFields(bbox);
            }
        }
    },

    setResultOption:function (resultOpt){
        this.resultOption = resultOpt;
        if(this.resultOption =='query'){
            jQuery('#submitQuery').removeClass('hide').addClass('inline-block');
            jQuery('#queryOutputFormat').removeClass('hide').addClass('inline-block');
            jQuery('#outputFormat').trigger('change');//trigger change to show the image options as needed
            jQuery('#submitMap').removeClass('inline-block').addClass('hide');
            //clear map
            qsInteractive.resetMap();
        }else{//map
            jQuery('#submitQuery').removeClass('inline-block').addClass('hide');
            jQuery('#queryOutputFormat').removeClass('inline-block').addClass('hide');
            jQuery('#image-options').removeClass('inline-block').addClass('hide');
            jQuery('#submitMap').removeClass('hide').addClass('inline-block');
            qsInteractive.hideQueryResult();
        }
    },

    setLocationOption:function (locOpt,fill){
        var locArray = null;
        this.selectedRegion = locOpt;
        this.locationByMap = false;
        if(locOpt){
            if(locOpt == "mapextent"){//map extent
                this.locationByMap = true;
                if(quakesMapApp.map){
                    locArray = quakesMapApp.map.getBBoxArray();
                }
            }else if(locOpt == "custom"){//custom
                fill = false;
                jQuery('#divLatLon').collapse("show");//show fields to enter
            }else{//regions
                var center = this.regionCenters[locOpt];
                if(center){
                    quakesMapApp.map.setView(center, 7);
                }
                var poly = this.regionPolies[locOpt];
                quakesMapApp.addPolygonlayer(poly);
            }
        }
        if(fill){
            this.fillLocationFields(locArray);
        }
    },
    checkLocationOption:function (locOpt){
        if(locOpt && locOpt == "mapextent"){//map extent
            this.locationByMap = true;
        }else{
            this.locationByMap = false;
        }
    },

    fillLocationFields:function (locArray){
        if(locArray && locArray.length ==4 ){
            jQuery("#minLon").val(this.padNumber(locArray[0],5));
            jQuery("#minLat").val(this.padNumber(locArray[1],5));
            jQuery("#maxLon").val(this.padNumber(locArray[2],5));
            jQuery("#maxLat").val(this.padNumber(locArray[3],5));
        }else{
            jQuery("#minLon").val('');
            jQuery("#minLat").val('');
            jQuery("#maxLon").val('');
            jQuery("#maxLat").val('');
        }
    },

    clearCoordinateFields:function(){
        this.fillLocationFields(null);
    },

    padNumber:function(numStr, numDigit){
        if(numStr){
            return new Number(numStr).toFixed(numDigit);
        }else{
            return '';
        }
    },

    fillDateFields:function (datesOpt){
        var endDate = new Date();
        //set to next hour if it passed current hour
        if(endDate.getUTCMinutes() > 0){
            endDate.setUTCHours(endDate.getUTCHours() + 1);
        }
        var startDate =  new Date();

        if(datesOpt){
            if(datesOpt == "a"){//last week
                startDate.setDate(startDate.getDate() - 7);
            }else if(datesOpt == "c"){//last year
                startDate.setFullYear(startDate.getFullYear() - 1);
            }else{
                startDate.setMonth(startDate.getMonth() - 1); //set to last month
            }
        }else{
            startDate.setMonth(startDate.getMonth() - 1); //set to last month
        }
        var startDateStr = startDate.getUTCFullYear()  + '-'+ (1 + startDate.getUTCMonth()) + '-' + startDate.getUTCDate();
        var endDateStr = endDate.getUTCFullYear()  + '-'+ (1 + endDate.getUTCMonth()) + '-' + endDate.getUTCDate()
        jQuery("#dateTo").val(endDateStr);
        jQuery("#dateFrom").val(startDateStr);
        jQuery("#hourTo").val(endDate.getUTCHours());
        jQuery("#hourFrom").val(startDate.getUTCHours());
        //update datepicker
        jQuery('#dateFrom').parent().data({
            date: startDateStr
        });
        jQuery('#dateTo').parent().data({
            date: endDateStr
        });
        jQuery('.date').datepicker('update');

    },

    /** clearDateFields */
    clearDateFields:function  (){
        jQuery("#dateTo").val('');
        jQuery("#dateFrom").val('');
        jQuery("#hourTo").val('');
        jQuery("#hourFrom").val('');
        jQuery("input[name='dateRadios']").each(function(){
            jQuery(this).attr("checked" , false );
        });
        this.resetMap();
    },
    /** clearDateFields */
    clearLocationFields:function(fill){
        this.setLocationOption(false,fill);
        jQuery("input[name='locRadios']").each(function(){
            jQuery(this).attr("checked" , false );
        });
        //reset map
        this.resetMap();
        quakesMapApp.map.setView(this.mapCentreNZ, 5);
    },
    /** reset map */
    resetMap:function(){
        //clear overlays
        var locOpt = jQuery("input[name='locRadios']:checked").val();
        var clearRegion = (!locOpt || 'mapextent' == locOpt);
        quakesMapApp.clearOverLays(clearRegion);
        //clear legend and result count
        jQuery('#mapLegend').removeClass('inline-block').addClass('hide');
        jQuery("#result-count").html("");
        this.showMapShow();
    },

    clearMagDepthFields:function (){
        jQuery("#minMag").val('');
        jQuery("#maxMag").val('');
        jQuery("#minDepth").val('');
        jQuery("#maxDepth").val('');
        this.resetMap();
    },

    clearQualityFields:function (){
        jQuery("#maxStdError").val('');
        jQuery("#minNumPhases").val('');
        jQuery("#minNumStations").val('');
    },

    initSelectionBoxes:function() {
        //2. init formats
        var listFormatOptions= "";
        for(var opt in this.outputOptions){
            listFormatOptions+= "<option value='" + this.outputOptions[opt] + "'>" + this.outputOptions[opt] + "</option>";
        }
        jQuery("#outputFormat").html(listFormatOptions);

        //imageSizes
        var listImageSizes= "";
        for(var opt in this.imageSizes){
            listImageSizes+= "<option value='" + this.imageSizes[opt] + "'>" + this.imageSizes[opt] + "px</option>";
        }
        jQuery("#imgSizeSelect").html(listImageSizes);
        //maxQuakesSelect
        var listMaxQuakes= "";
        for(var opt in this.maxNumQuakes){
            listMaxQuakes+= "<option value='" + this.maxNumQuakes[opt] + "'>" + this.maxNumQuakes[opt] + "</option>";
        }
        jQuery("#maxQuakesSelect").html(listMaxQuakes);

        //hourFrom hourTo
        var listHour = "<option value='" + "'></option>";
        for(var i = 0; i < 24; i++){
            var zeropaded = i;
            if(i < 10){
                zeropaded = "0" + i;
            }
            listHour+= "<option value='" + i + "'>" + zeropaded + ":00</option>";
        }
        jQuery("#hourFrom").html(listHour);
        jQuery("#hourTo").html(listHour);
    },

    makeSearchQuery:function(){
        var loc = window.location;
        qsInteractive.currentURL = loc.protocol + '//' + loc.host + loc.pathname.substr(0, loc.pathname.lastIndexOf("/")) + "/";
        var lastIndex = qsInteractive.currentURL.lastIndexOf("/");
        if(lastIndex < (qsInteractive.currentURL.length - 1)){
            qsInteractive.currentURL = qsInteractive.currentURL.substr(0, lastIndex + 1);
        }
        //console.log("loc.pathname " + loc.pathname.substr(0, loc.pathname.lastIndexOf("/")));
        qsInteractive.queryString = "", qsInteractive.queryStringDate = '';

        //1. date
        var dateStart = jQuery.trim(jQuery("#dateFrom").val());
        if(dateStart != ''){
            var hourStart = jQuery.trim(jQuery("#hourFrom").val());
            if(hourStart != ''){
                dateStart += "T" + hourStart + ":00:00";
            }
            qsInteractive.queryStringDate = qsInteractive.appendQueryString(qsInteractive.queryStringDate, "startdate=" + dateStart);
        }
        var dateEnd = jQuery.trim(jQuery("#dateTo").val());

        if(dateEnd != ''){
            var hourEnd = jQuery.trim(jQuery("#hourTo").val());
            if(hourEnd != ''){
                dateEnd += "T" + hourEnd + ":00:00";
            }
            qsInteractive.queryStringDate = qsInteractive.appendQueryString(qsInteractive.queryStringDate, "enddate=" + dateEnd);
        }
        //2. region/bbox
        if(this.selectedRegion){
            if(this.selectedRegion =='mapextent' || this.selectedRegion =='custom' ){
                var minLat = jQuery.trim(jQuery("#minLat").val());
                var minLon = jQuery.trim(jQuery("#minLon").val());
                var maxLat = jQuery.trim(jQuery("#maxLat").val());
                var maxLon = jQuery.trim(jQuery("#maxLon").val());

                if(minLat != '' && minLon != '' && maxLat != '' && maxLon != ''){
                    var bbox = minLon + ',' + minLat + ',' + maxLon + ',' + maxLat ;
                    qsInteractive.queryString = qsInteractive.appendQueryString(qsInteractive.queryString, "bbox=" + bbox);

                }
            }else{//region
                qsInteractive.queryString = qsInteractive.appendQueryString(qsInteractive.queryString, "region=" + this.selectedRegion);
            }
        }

        //3. magnitude
        var minMag = jQuery.trim(jQuery("#minMag").val());
        var maxMag = jQuery.trim(jQuery("#maxMag").val());
        if(minMag != '') qsInteractive.queryString = qsInteractive.appendQueryString(qsInteractive.queryString, "minmag=" + minMag);
        if(maxMag != '') qsInteractive.queryString = qsInteractive.appendQueryString(qsInteractive.queryString, "maxmag=" + maxMag);

        //4. depth
        var maxDepth = jQuery.trim(jQuery("#maxDepth").val());
        var minDepth = jQuery.trim(jQuery("#minDepth").val());
        if(minDepth != '') qsInteractive.queryString = qsInteractive.appendQueryString(qsInteractive.queryString, "mindepth=" + minDepth);
        if(maxDepth != '') qsInteractive.queryString = qsInteractive.appendQueryString(qsInteractive.queryString, "maxdepth=" + maxDepth);
        //get result format
        qsInteractive.resultFormat = jQuery("#outputFormat").val().toLowerCase();

    },

    /*** show search result for map ***/
    getMapSearchResult:function(){
        this.makeSearchQuery();
        var queryString = qsInteractive.appendQueryString(qsInteractive.queryString, qsInteractive.queryStringDate);
        //console.log(  " queryString " + queryString );
        //5. max quakes maxQuakes
        var maxQuakes = jQuery.trim(jQuery("#maxQuakesSelect").val());
        if(maxQuakes ){
            maxQuakes = new Number(maxQuakes);
        }
        qsInteractive.resultMaxQuakes = maxQuakes;
        qsInteractive.resultDataUrl = qsInteractive.currentURL + qsInteractive.SERVICE_PATH + 'geojson?limit=' + maxQuakes + '&' + queryString;
        qsInteractive.resultCountUrl = qsInteractive.currentURL + qsInteractive.SERVICE_PATH + 'count?' + queryString;
        //
        qsInteractive.queryQuakesCount(qsInteractive.resultCountUrl, 1);
    },

    /*** show search result for query ***/
    getQuerySearchResult:function(){
        this.makeSearchQuery();
        var queryString = qsInteractive.appendQueryString(qsInteractive.queryString, qsInteractive.queryStringDate);

        qsInteractive.resultDataUrl = qsInteractive.currentURL + qsInteractive.SERVICE_PATH + qsInteractive.getPathFormatOption() + '?' + queryString;
        qsInteractive.resultCountUrl = qsInteractive.currentURL + qsInteractive.SERVICE_PATH + 'count?' + queryString;

        //get quakes count
        qsInteractive.queryQuakesCount(qsInteractive.resultCountUrl, 0);

    },

    getPathFormatOption:function(){
        var path = this.resultFormat;
        if('image' == this.resultFormat){
            path += '/' + this.imageSize;
            this.mapShowPlaces = jQuery("#mapTownCheck")[0].checked;
            if(this.mapShowPlaces){
                path += '/1'
            }else{
                path += '/0'
            }
        }
        return path;
    },

    appendQueryString:function(queryString, value){
        if(queryString ==''){
            return value;
        }else{
            if(value != ''){
                queryString += '&' + value;
            }
            return queryString;
        }
    },

    getQuakesCountLabel:function(count, styleclass){
        if(!styleclass){
            styleclass = count;
        }
        if(styleclass == 0){//no data
            return "<span class='badge'>" + count + "</span>";
        }else if(styleclass < 20000){//ok
            return "<span class='badge badge-success'>" + count + "</span>";
        }else if(styleclass < 10000){//warning
            return "<span class='badge badge-warning'>" + count + "</span>";
        }else{//severe
            return "<span class='badge badge-important'>" + count + "</span>";
        }
    },

    showMapResult:function(result){
        //console.log("result " + result)
        var quakesCount = result.count;
        var rows = 2;
        if(quakesCount != null && quakesCount != "undefined"){//0 is ok
            var resultCountString = "Returned Quakes: " + this.getQuakesCountLabel(quakesCount, 100);
            if(qsInteractive.resultMaxQuakes && quakesCount > 0){
                var displayedQuakes = qsInteractive.resultMaxQuakes;
                if(displayedQuakes > quakesCount){
                    displayedQuakes = quakesCount;
                }
                resultCountString = "Showing " +  this.getQuakesCountLabel(displayedQuakes, 100) + " of " + this.getQuakesCountLabel(quakesCount, 100) + " returned quakes.";
            }
            jQuery("#result-count").html(resultCountString);
        }
        //
        if(quakesCount > 0){
            quakesMapApp.loadQuakesData(qsInteractive.resultDataUrl);
        }
    },

    showQueryResult:function(result){
        //console.log("result " + result)
        var quakesCount = result.count;
        var rows = 2;
        if(quakesCount != null && quakesCount != "undefined"){//0 is ok
            jQuery("#result-number").html(this.getQuakesCountLabel(quakesCount) + " Returned Quakes");
        }
        var resultUrl = '#', resultUrlAll = '';
        jQuery("#result-details").css("display", "block");
        if(quakesCount == 0){//hide result-details panel
            jQuery("#result-details").css("display", "none");
        }else if(quakesCount <= qsInteractive.MAX_NUMBER_QUAKES){
            resultUrl = qsInteractive.resultDataUrl;
            resultUrlAll = qsInteractive.resultDataUrl;
        }else{
            if(result.dates && result.dates.length > 1){
                var start = result.dates[1], end = result.dates[0];
                var queryStrDate;
                // take the first one for downloading
                resultUrl = qsInteractive.currentURL + qsInteractive.SERVICE_PATH + qsInteractive.getPathFormatOption() + '?'
                + qsInteractive.appendQueryString(qsInteractive.queryString, "startdate=" + start + "&enddate=" + end) ;
                for(var i = 1; i < result.dates.length; i++){
                    start = result.dates[i];
                    //make multiple request urls
                    resultUrlAll += qsInteractive.currentURL + qsInteractive.SERVICE_PATH + qsInteractive.getPathFormatOption() + '?'
                    + qsInteractive.appendQueryString(qsInteractive.queryString, "startdate=" + start + "&enddate=" + end) + "\n";
                    end = start;
                    rows++;
                }
            }
        }
        jQuery("#result-text").val(resultUrlAll);
        jQuery("#btnDownload").attr("href", resultUrl);
        jQuery("#result-text").attr("rows", rows);

        jQuery("#result").css("display", "inline-block");
        jQuery("#btnReset").css("display", "inline-block");
        if(rows > 2){
            jQuery("#result p.text-warning").css("display", "inline-block");
        }
        else{
            jQuery("#result p.text-warning").css("display", "none");
        }
        if(!quakesCount){
            jQuery("#btnDownload").removeClass("active").addClass("hide");
        }else{
            jQuery("#btnDownload").removeClass("hide").addClass("active");
        }
    },

    resetFields:function(start){
        jQuery("#result textarea").val("");
        jQuery("input:radio[name='dateRadios']")[1].checked = true;
        qsInteractive.fillDateFields("b");//last month
        qsInteractive.clearMagDepthFields();
        qsInteractive.clearCoordinateFields();
        jQuery("input:radio[name='dateRadios']")[1].checked = true;
        jQuery("input:radio[name='locRadios']")[0].checked = true;
        if(start) {
            jQuery("input:radio[name='resultRadios']")[0].checked = true;
        }
        jQuery("#mapTownCheck").checked = false;
        //qsInteractive.setLocationOption('m',true);
        this.hideQueryResult();
    },
    /*** change btnShowMap to clear ***/
    showMapClear:function(){
        jQuery('#btnShowMap').val('Clear');
        qsInteractive.locationByMap = false;
        jQuery('#mapLegend').removeClass('hide').addClass('inline-block');
    },
    /*** change btnShowMap to Show ***/
    showMapShow:function(){
        //jQuery("input[name='locRadios']").trigger('change');//trigger change
        jQuery('#btnShowMap').val('Show');
        jQuery('#mapLegend').removeClass('inline-block').addClass('hide');
        //enable coordinates update
        qsInteractive.checkLocationOption(jQuery("input[name='locRadios']:checked").val());
    },

    /* check a number is entered, for year value */
    checkNum:function (e, allowVals) {
        e = e||window.event;
        var charCode = (navigator.appName == "Netscape") ? e.which : e.keyCode;
        //console.log("## charCode=" + charCode);
        if(allowVals){//allowVals
            for(var i in allowVals){
                if(charCode == allowVals[i]){
                    return true;
                }
            }
        }

        if (charCode > 31 && (charCode < 48 || charCode > 57)) {
            return false;
        }
        return true;
    },
    queryQuakesCount:function(url, resultopt){//
        //console.log("queryQuakesCount " );
        jQuery.getJSON(url, function(data) {
            //console.log("data " );
            if(resultopt == 0){//get query
                qsInteractive.showQueryResult(data);
            }else{//show on map
                qsInteractive.showMapResult(data);
            }
        });
    }
};

/**
 * define basic variables and functions for the application
 */
var quakesMapApp = {
    URL_PREFIX : "http://www.geonet.org.nz/quakes/region/newzealand/",
    GEOJSON_URL : "/quakes/services/quakes/newzealand",
    BING_KEY : '###',

    ieVersion:-1,
    currentTime:null,
    MAGCLASSES : ["Mag <2","Mag 2-3","Mag 3-4", "Mag 4-5","Mag 5-6","Mag 6-7","Mag \u22657"],//>=7
    maxQuakesList : ["100","500","1000","1500"],
    minIntensity : 3,
    maxQuakes:100,
    map: null,
    quakeStartTime:null,
    allQuakesLayers:[],
    regionLayer:null,
    layerControl:null,
    regionMap:false,

   //init map
   initMap:function(key) {
        if (key) {
            quakesMapApp.BING_KEY = key;
        }
        quakesMapApp.ieVersion = quakesMapApp.getIEVersion();

        quakesMapApp.map = new L.Map('haz-map', {
            attributionControl: false,
            worldCopyJump: false
        }).setView(qsInteractive.mapCentreNZ, 5);

        var osmGeonetUrl = 'http://{s}.geonet.org.nz/osm/tiles/{z}/{x}/{y}.png', //
                osmMqUrl = 'http://{s}.mqcdn.com/tiles/1.0.0/map/{z}/{x}/{y}.jpg', //
                cloudmadeAttribution = '',
                osmLayerGeonet = new L.TileLayer(osmGeonetUrl, {
                    minZoom: 1,
                    maxZoom: 18,
                    attribution: cloudmadeAttribution,
                    errorTileUrl: 'http://static.geonet.org.nz/osm/images/logo_geonet.png',
                    subdomains: ['static1', 'static2', 'static3', 'static4', 'static5']
                });

        var bingLayer = new L.BingLayer(quakesMapApp.BING_KEY, {type: "Aerial"});

        //map switcher
        var baseLayers = {
            "Map": osmLayerGeonet,
            "Aerial": bingLayer
        };

        quakesMapApp.map.addLayer(osmLayerGeonet);
        //add layer switch
        quakesMapApp.layerControl = L.control.layers(baseLayers);
        // force refresh ovverlays on ie
        quakesMapApp.layerControl.refreshOverlay = function () {
            var i, input, obj,
                inputs = this._form.getElementsByTagName('input'),
                inputsLen = inputs.length;
            this._handlingClick = true;
            for (i = 0; i < inputsLen; i++) {
                input = inputs[i];
                obj = this._layers[input.layerId];
                if (input.checked && this._map.hasLayer(obj.layer)) {
                    if (obj.overlay) {//turn off and then turn on the overlay
                        this._map.removeLayer(obj.layer);
                        this._map.addLayer(obj.layer);
                    }
                }
            }
            this._handlingClick = false;
        };
        // force expand layer control
        quakesMapApp.layerControl.changeLayout = function (expand) {
            this.options.collapsed = !expand;
            if (expand) {//remove listener
                L.DomEvent.off(this._container, 'mouseout', this._collapse);
                this._map.off('movestart', this._collapse);
                this._expand();
            } else {//add listener
                L.DomEvent.on(this._container, 'mouseout', this._collapse, this);
                this._map.on('movestart', this._collapse, this);
                this._collapse();
            }
        };

        quakesMapApp.layerControl.addTo(quakesMapApp.map)
        //check coordinates
        quakesMapApp.map.getBBoxArray = function () {
            var bound = this.getBounds();
            if (bound) {
                return bound.toBBoxString().split(',')
            }
            return null;
        };

        quakesMapApp.map.on('moveend', function (e) {
            quakesMapApp.updateMapBound();
        });

    },

    clearOverLays:function(clearRegion){
        if(this.allQuakesLayers.length > 0){
            for(var index in this.allQuakesLayers){
                this.map.removeLayer(this.allQuakesLayers[index]);
                this.layerControl.removeLayer(this.allQuakesLayers[index]);
                this.allQuakesLayers[index].clearLayers();
            }
            this.allQuakesLayers = [];
        }
        if(clearRegion && this.regionLayer){
            this.map.removeLayer(this.regionLayer);
        }
    },

    checkMapData:function(minIntensitySel, maxQuakesSel){
        var update = false;
        if(minIntensitySel != this.minIntensity){
            this.minIntensity = minIntensitySel;
            update = true;
        }
        if(maxQuakesSel != this.maxQuakes){
            this.maxQuakes = maxQuakesSel;
            update = true;
        }
        if(update){
            this.loadQuakesData();
        }
    },

    loadQuakesData:function(url){
        var count = 0;
        jQuery.getJSON(url, function (data) {
            var dataByMagClass = {};
            //classify geojson data for each intensities
            for(var num = 0; num < data.features.length; num++){
                var feature = data.features[num];
                var magClass = quakesMapApp.getQuakeMagClass(1*feature.properties.magnitude);
                if(!dataByMagClass[magClass]){
                    dataByMagClass[magClass] = {
                        "type": "FeatureCollection",
                        "features": []
                    };
                }
                dataByMagClass[magClass].features.push(feature);
                //console.log("##0 magClass " + magClass);
                count++;
            }
            //make quakes layers for each magClass level, and add to maps
            quakesMapApp.clearOverLays(false);
            for(var index in quakesMapApp.MAGCLASSES){
                var magClass = quakesMapApp.MAGCLASSES[index];
                if(dataByMagClass[magClass]){
                    var layer = quakesMapApp.makeQuakesLayer(dataByMagClass[magClass]);
                    quakesMapApp.map.addLayer(layer);
                    var layerName =  magClass;
                    quakesMapApp.layerControl.addOverlay(layer, layerName);
                    quakesMapApp.allQuakesLayers.push(layer);
                }
            }
            //show clear button
            qsInteractive.showMapClear();
        });
    },
    addPolygonlayer:function(polyCoords){
        if(this.regionLayer){
            quakesMapApp.map.removeLayer(this.regionLayer);
        }
        if(polyCoords){//add region
            this.regionLayer =  L.polygon(polyCoords, {
                color: '#D0D0D0',
                clickable: false,
                fillOpacity: 0.1
            }).addTo(quakesMapApp.map);
        }else{
            this.regionLayer = null;
        }
        return this.regionLayer;
    },
    makeQuakesLayer:function(data){
        return L.geoJson(data, {
            onEachFeature: function (feature, layer) {
                if (feature.properties) {
                    var quaketime = Date.parseUTC(feature.properties.origintime);
                    var nzTime = Date.toNZTimeString(quaketime);
                    var popup = "<strong>Public ID: </strong><a target='_blank' href='" + quakesMapApp.URL_PREFIX + feature.properties.publicid + "'>" + feature.properties.publicid + "</a><br/>" +
                    "<strong>Time: </strong>" + nzTime + "<br/>" +
                    "<strong>Magnitude: </strong>" + qsInteractive.padNumber(feature.properties.magnitude,1) + "<br/>" +
                    "<strong>Depth: </strong>" + qsInteractive.padNumber(feature.properties.depth,0) + " km<br/>";
                    layer.bindPopup(popup);
                }
            },
            pointToLayer: function (feature, latlng) {
                //fix dateline issue
                var lon;
                if (latlng.lng < 0) {
                    lon = latlng.lng + 360;
                }else {
                    lon = latlng.lng;
                }
                var newLatLong = new L.LatLng(latlng.lat, lon, true);
                var quakecolor = quakesMapApp.getQuakeColor(1*feature.properties.depth);
                var quakeSize = quakesMapApp.getQuakeSize(1*feature.properties.magnitude);
                if(!quakecolor){
                    quakecolor = '#F5F5F5';
                }
                if(!quakeSize){
                    quakeSize = 5;
                }

                return L.circleMarker(newLatLong, {
                    radius: quakeSize,
                    color: quakecolor,
                    fillColor: quakecolor,
                    stroke: false,
                    zIndexOffset: 0.1*(10000 - 1*feature.properties.depth),
                    fillOpacity: 0.7,
                    opacity: 0.7
                });
            }
        });
    },

    updateMapBound:function(){
        if(this.map){
            var bound = this.map.getBBoxArray();
            if(bound){
                qsInteractive.fillMapBound(bound);
            }
        }
    },

    MAGS4CLASS : [2, 3, 4, 5, 6, 7],
    getQuakeMagClass:function(mag){
        if(!mag){//dodgy value for unnoticeable quakes
            mag = -9;
        }
        for (var i = 0; i < this.MAGS4CLASS.length; i++) {
            if (mag < this.MAGS4CLASS[i]) {
                return this.MAGCLASSES[i];
            }
        }
        return this.MAGCLASSES[this.MAGS4CLASS.length];
    },

    MAGS4SIZE : [2, 3, 4, 5, 6, 7],
    MAGSIZES : [3, 5, 7, 9, 11, 13, 15],
    getQuakeSize:function(mag){
        if(!mag){//dodgy value for unnoticeable quakes
            mag = -9;
        }
        for (var i = 0; i < this.MAGS4SIZE.length; i++) {
            if (mag < this.MAGS4SIZE[i]) {
                return this.MAGSIZES[i];
            }
        }
        return this.MAGSIZES[this.MAGS4SIZE.length];
    },

    DEPTHS : [15, 40, 100, 200],
    QUAKECOLORS : ['#B5FCFF', '#70F0FF', '#00B9CE', '#008291','#004951'],
    getQuakeColor:function(depth){
        for (var i = 0; i < this.DEPTHS.length; i++) {
            if (depth < this.DEPTHS[i]) {
                return this.QUAKECOLORS[i];
            }
        }
        return this.QUAKECOLORS[this.DEPTHS.length];
    },

    // Returns the version of Windows Internet Explorer or a -1
    // (indicating the use of another browser).
    getIEVersion:function(){
        var rv = - 1;
        // Return value assumes failure.
        if (navigator.appName == 'Microsoft Internet Explorer') {
            var ua = navigator.userAgent;
            var re = new RegExp("MSIE ([0-9]{1,}[\.0-9]{0,})");
            if (re.exec(ua) != null)
                rv = parseFloat(RegExp.$1);
        }
        return rv;
    }
};


/*** app starts heres !!! ***/
function initPageElements(key) {
    qsInteractive.initDatePicker();
    qsInteractive.initSelectionBoxes();
    qsInteractive.initEvents();
    //init map
    quakesMapApp.initMap(key);
    qsInteractive.resetFields(true);
    //populate coordinates fields with map extents
    qsInteractive.setLocationOption('mapextent',true);
    //IE right issue
    if(quakesMapApp.ieVersion > -1){
        jQuery("#btnShowMap").removeClass("pull-right");
        //IE date picker width issue
        jQuery(".datepicker table ").css("width", "150px");
    }
    // init modal views
    jQuery("#myModal").on('hidden', function () {
        jQuery("#haz-map").insertBefore("#mapLegend");
        //jQuery("#mapContainer br").insertBefore("#haz-map");
        jQuery("#haz-map").removeClass("large-quake-search-map").addClass("small-quake-search-map");

        if(quakesMapApp.ieVersion > -1 && quakesMapApp.ieVersion < 9){
            quakesMapApp.layerControl.refreshOverlay();
        }
        quakesMapApp.layerControl.changeLayout(false)
        quakesMapApp.map.invalidateSize();
    })

    jQuery("#myModal").on('shown', function () {
        jQuery("#myModal .modal-header").css('background-color', '#F5F5F5');
        jQuery("#haz-map").insertBefore("#myModal .modal-footer");
        jQuery("#haz-map").removeClass("small-quake-search-map").addClass("large-quake-search-map");
        if(quakesMapApp.ieVersion > -1 && quakesMapApp.ieVersion < 9){
            quakesMapApp.layerControl.refreshOverlay();
        }
        quakesMapApp.layerControl.changeLayout(true)
        quakesMapApp.map.invalidateSize();
        var isMobile = navigator.userAgent.match(/(iPhone|iPod|iPad|Android|BlackBerry)/);
        if (!isMobile) {
            //scroll to top otherwise it offsets the zoom position for the map
            //don't scroll for tablets since it causes some rendering issues
            scroll(0, 0);
        }
    })
}

function closeLargeMap() {
    jQuery("#myModal").modal('hide');
}

function showLargeMap() {
    jQuery("#myModal").modal('show');
}


