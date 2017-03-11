
library routes;

import 'dart:html';
import 'package:route_hierarchical/client.dart';
import 'userAuth.dart';
import 'formHandler.dart';
import 'photoUpload.dart';

var panelCount = 1;


// Setup routing for pretty urls.
setupRouters() async {
  var router = new Router();
  router.root
    ..addRoute(name: 'newAddress', path:'/newAddress', enter: newAddressRoute)
    ..addRoute(name: 'userNotFound', path:'/userNotFound', enter: userNotFoundRoute)
    ..addRoute(name: 'validationFailed', path:'/validationFailed', enter: setRoute)
    ..addRoute(name: 'about', path: '/about', enter: setRoute)
    ..addRoute(name: 'journal', path: '/journal', enter: setRoute)
    ..addRoute(name: 'reports', path: '/reports', enter: setRoute)
    ..addRoute(name: 'resendValidation', path: '/resendValidation', enter: setRoute)
    ..addRoute(name: 'home', path: '/', defaultRoute: true, enter: setRoute)
    ..addRoute(name: 'signOut', path: '/signOut', enter: signOut);
  router.listen();
}

userNotFoundRoute(RouteEvent e) {
   var signinLink = querySelector('#signin');
   signinLink.click();
   var signinWarn = querySelector('#signinWarn');
   signinWarn.innerHtml = "Not able to sign you in automatically, please sign in to continue.";
   signinWarn.style.display = "block";
   ModElement signinModal = querySelector('#signinModal');

}

newAddressRoute(RouteEvent e) {
  var signinLink = querySelector('#signin');
  signinLink.click();
  var signinWarn = querySelector('#signinWarn');
  signinWarn.innerHtml = "Visiting the site from a new location, please sign in to continue.";
  signinWarn.style.display = "block";
  ModElement signinModal = querySelector('#signinModal');
}

addJournalListeners() async {
  panelCount = 1;
  var trophyListener = querySelector('#add-trophy-button');
  trophyListener.onClick.listen((e) {
    addPanel('trophy-panel', 'views/journal-trophy.html');
  }); 
  var fishListener = querySelector('#add-fish-button');
  fishListener.onClick.listen((e) {
    addPanel('fish-panel', 'views/journal-fish.html');

  });
  var journalSubmitListener = querySelector('#journal_submit');
  journalSubmitListener.onClick.listen((e){
    formSubmit(querySelector("#journalEntry"));
  });
}

// Call sets body div to html file for routing.
setRoute(RouteEvent e) {
  DivElement body = querySelector('#body');

  HttpRequest.getString('/views/' + e.route.name + '.html').then((bodyHTML){
    body.setInnerHtml(bodyHTML, treeSanitizer: NodeTreeSanitizer.trusted);
  }).then((t){
    // Add listeners for journal entry page
    if (e.route.name == 'journal') {
      addJournalListeners();
    }
    populateSelects();
  });
}


// Onclick to add form to panel from HTML
addPanel(String panelId, String page) async {
  // Add new element to panel.
  Element panelDiv = querySelector('#${panelId}');
  String panelHtml = await HttpRequest.getString(page);
  await panelDiv.appendHtml(panelHtml, treeSanitizer: NodeTreeSanitizer.trusted);

  Element newPanel = panelDiv.querySelector('#new-panel');

  // Set incremented ID for panel
  newPanel.id = panelDiv.id  + "_" + panelCount.toString();

  window.console.debug(newPanel.id);

  //Update names of the input fields
  ElementList inputs = newPanel.querySelectorAll('[name]');
  for (Element input in inputs) {
    input.id = input.id + "_" + panelCount.toString();
    window.console.debug(input.id);
  }


  // Add listener to remove panel button
  Element newRemove = newPanel.querySelector('#remove-panel');
  newRemove.id = newRemove.id + "_" + panelDiv.id + "_" + (panelCount++).toString();

  window.console.debug(newRemove.id);
  newRemove.onClick.listen((e){
    newPanel.remove();
  });
  populateSelects();

  // Add upload photo functionality
  if(panelId == 'trophy-panel'){
    addPhotoListener(newPanel.id);
  }

}