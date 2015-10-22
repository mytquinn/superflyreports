import 'dart:html';
import 'package:bootjack/bootjack.dart';

final navLinks = {
	'home':'views/home.html',
	'about':'views/about.html',
	'journal':'views/journal.html',
	'reports':'views/reports.html'
	
};

main() async {
   Modal.use();
   Transition.use();
   DivElement header = querySelector('#header');
   String headerHTML = await HttpRequest.getString("views/header.html");
   header.setInnerHtml(headerHTML, treeSanitizer: NodeTreeSanitizer.trusted);

   DivElement footer = querySelector('#footer');
   String footerHTML = await HttpRequest.getString("views/footer.html");
   footer.setInnerHtml(footerHTML, treeSanitizer: NodeTreeSanitizer.trusted);

   requestBody('views/about.html');
   
   addModal('signup','signupModal','views/signup.html');
   addModal('signin','signinModal','views/signin.html');

   var btnListeners = querySelectorAll('.navbar-btn');
   for (var btnListener in btnListeners) {
      btnListener.onClick.listen((e) {
	 requestBody(navLinks[btnListener.id]);
      });
   }
}

requestBody(String page) async {
   DivElement body = querySelector('#body');
   String bodyHTML = await HttpRequest.getString(page);
   body.setInnerHtml(bodyHTML, treeSanitizer: NodeTreeSanitizer.trusted);
}

addModal(String modalLinkId, String modalId, String page) async {
	 Element modalDiv = querySelector('#modal');
         Element modalLink = querySelector('#' + modalLinkId);

	 modalLink.attributes['data-toggle'] = 'modal';
	 modalLink.attributes['data-target'] = '#' + modalId;
	 String modalHtml = await HttpRequest.getString(page);
	 modalDiv.appendHtml(modalHtml, treeSanitizer: NodeTreeSanitizer.trusted);
}


