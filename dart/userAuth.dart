library userAuth;

import 'package:http/browser_client.dart';
//import 'package:http/http.dart' as http;
import 'dart:html';
import 'package:route_hierarchical/client.dart';

var sessionId = null;

bool checkLogin() async {
  var client = new BrowserClient();
  if (!document.cookie.contains('userSession_id')) {
       return false;
  }
  var value;
  var url = '/checkLogin';
  await client.get(url).then((response) {
    value = response.body;
  });
  if (value == 'true') {
    return true;
  } else {
    return false;
  }
}

autoSignin() async {
  var client = new BrowserClient();
  var url = '/doSignin';
  await client.get(url);
}

signOut(RouteEvent e) async {
  var client = new BrowserClient();
  var url = '/doSignout';
  await client.get(url);
  document.cookie = '';
  window.location.href = '/';
}

getSessionData(String data) async {
  var client = new BrowserClient();
  var value;
  var url = '/getSessionData/${data}';
  await client.get(url).then((response) {
    value = response.body;
  });
  return value;
}
