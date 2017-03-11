library formHandler;

import 'package:http/browser_client.dart';
import 'dart:html';
import 'dart:convert';
import 'dart:core';

populateSelects() async {
  var selects = querySelectorAll('select');
  for (var select in selects) {
    if (select.attributes.containsKey('datasrc') && select.children.length <= 1) {
      var table = select.attributes['datasrc'].split('.')[0];
      var field = select.attributes['datasrc'].split('.')[1];
      var client = new BrowserClient();
      var url = '/getType/${table}';
      var jsonResponse;
      await client.get(url).then((response){
        Map parsedJson = JSON.decode(response.body);
        for (var i=0; i < parsedJson.length; i++){
          select.appendHtml('<option value=${parsedJson[i]["id"]}>${parsedJson[i][field]}</option>');
        }
      });
    }
  }
}

formSubmit (FormElement form) {
  window.console.debug(form);

  ElementList inputs = form.querySelectorAll('[datafld]');
  window.console.debug(inputs);
  int tableCount = 0;
  String field ='';
  String value = '';
  String table = '';
  bool newTable = false;
  Map formMap = new Map();
  Map tableMap = new Map();
  for (Element input in inputs) {
    String datafld = input.attributes['datafld'];
    if(datafld.contains('table:')) {
      datafld = datafld.replaceFirst('table:', '');
      table = datafld.split('.')[0];
      if(datafld.contains('=')) {
        field = datafld.split('=')[0].split('.')[1];
        value = datafld.split('=')[1];
      }
      newTable = true;

    } else {
      table = datafld.split('.')[0];
      field = datafld.split('.')[1];
      if (input.tagName != 'LI') {
        value = input.innerHtml;
      } else {
        value = input.text;
      }
    }
    if(newTable) {
      tableMap.clear();
      tableCount++;
      newTable = false;

    } else {
       if (value == '') { value = 'null'; }
       tableMap[field] = value;
       window.console.debug(tableMap);
    }
  }
  formMap[table.toString()] = tableMap.toString();
  final JsonString =  JSON.encode(formMap).toString();
  final JsonFinal = JsonString.replaceAllMapped(new RegExp(r'-\d+'), (match) {
    return '';

  });
  window.console.debug(JsonFinal);
}




