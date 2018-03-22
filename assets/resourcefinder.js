var t, address;
var resources = {};
var system = require('system');

var getHostname = function(str) {
  // stollen from http://beardscratchers.com/journal/using-javascript-to-get-the-hostname-of-a-url
  try{
    var re = new RegExp('^(?:f|ht)tp(?:s)?\://([^/]+)', 'im');
    return str.match(re)[1].toString();
  } catch (err) {
    //probably data uri which we dont care abt
    //console.log(str)
    return null;
  }
}

//No need for headers for each individual object
//Simplify output
//Match output format of http.Headers
var reduceheaders = function(headers){
  var newheaders = {}
  var key, value, i;
  var cdn = null;
  //console.log(JSON.stringify(headers))
  for(i=0;i<headers.length;i++){
    key = headers[i].name;
    value = headers[i].value;
//       console.log(key, value)
    newheaders[key] = [value];
  }
  return newheaders;
}

var makereport = function(input){
  var key, keys, basepagedomain, output, i;
//    console.log("making report")
  keys = Object.keys(input.resources);
  //console.log(JSON.stringify(keys))
  for (i=0;i<keys.length;i++){
    key = keys[i];
    input.resources[key].headers = reduceheaders(input.resources[key].headers);
    input.resources[key].isbase = input.basepagehost == key;
    input.resources[key].hostname = key;
    //delete input.resources[key].headers;
  }
  //nodejs app reads from console
  console.log(JSON.stringify(input));
}


t = Date.now();
address = system.args[1];

var loadpage = function(url){
  var page = require('webpage').create()
  //https://github.com/ariya/phantomjs/issues/10389#issuecomment-103650123
  page.onNavigationRequested = function(url, type, willNavigate, main) {
      if (main && url!=address) {
          address = url;
          //console.log("redirect caught")
          page.close()
          setTimeout('loadpage(address)',1); //Note the setTimeout here
      }
  };

  page.onResourceReceived = function(request){
    var url, size, hostname, headers, i;
    var headers = request.headers;
  //        console.log(JSON.stringify(request));
    url = request.url;
    if (!(size)){
      size = (request.bodySize ? request.bodySize: 0);
    }
    //console.log(url);
    hostname = getHostname(url);
    if ((hostname) && (size > 0)){
      if (!(resources[hostname])){
        resources[hostname] = {};
        resources[hostname].count = 0;
        resources[hostname].bytes = 0;
      }
      resources[hostname].count += 1;
      //phantomjs lies! so we see content-length header is available
      for (i=0;i<headers.length;i++){
        if (headers[i].name.toLowerCase() == "content-length"){
          size = parseInt(headers[i].value);
          break;
        }
      }

      resources[hostname].bytes += size;
      //save the last response headers per host
      resources[hostname].headers = headers;
    }
  }

  //Silently ignore js err
  page.onError = function(msg, trace) {
  }

  page.open(address, function (status) {
    var output;
    if (status !== 'success') {
      console.log('{"error": "FAIL"}');
    } else {
      t = Date.now() - t;
  //            console.log('Loading time ' + t + ' msec');
      output = {};
      output.basepagehost = page.evaluate(function () {
          return document.location.hostname;
      });

      output.resources = resources;
      makereport(output);
      //console.log(JSON.stringify(output));
    }
    phantom.exit();
  });
}

loadpage(address);
