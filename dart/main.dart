import 'dart:html';
import 'package:bootjack/bootjack.dart';
import 'routes.dart';
import 'userAuth.dart';

final navLinks = {
  'home':'views/home.html',
  'about':'views/about.html',
  'journal':'views/journal.html',
  'reports':'views/reports.html'
};

main() async {

  await autoSignin();
  //Instantiate modules for bootjack
  Modal.use();
  Transition.use();
  Dropdown.use();
  Collapse.use();

  // Load header and footer on page
  DivElement header = querySelector('#header');
  String headerHTML = await HttpRequest.getString('views/header.html');
  await header.setInnerHtml(headerHTML, treeSanitizer: NodeTreeSanitizer.trusted);

  setupHeader();

  DivElement footer = querySelector('#footer');
  String footerHTML = await HttpRequest.getString('views/footer.html');
  await footer.setInnerHtml(footerHTML, treeSanitizer: NodeTreeSanitizer.trusted);
  await loadModals();
  setupRouters();
  configureSignin();

}

// Load Modals dialogs from html
loadModals() async{
  await addModal('signup','signupModal','views/signup.html');
  await addModal('signin','signinModal','views/signin.html');

}
// remove elements based on signed-in/signed-out classes
configureSignin() async{
  var loggedIn = await checkLogin();
  window.console.log(loggedIn);

  var toRemove;
  if (loggedIn) {
    toRemove = querySelectorAll('.signed-out');
    var username = await getSessionData('username');
    var usernameElement = querySelector('#username');
    usernameElement.innerHtml = username;
  } else {
    toRemove = querySelectorAll('.signed-in');
  }
  toRemove.forEach((element){
    element.remove();
  });



}
// Adds header behavior
setupHeader() {
  DivElement headerDiv = querySelector('#header-div');
  ImageElement titleImage = querySelector('#title-image');

  headerDiv.style.height = (window.innerWidth.toInt()/3).toString() + 'px';
  titleImage.style.width = (window.innerWidth.toInt()/3).toString() + 'px';
  titleImage.style.paddingTop = (window.innerWidth.toInt()/24).toString() + 'px';

  window.onResize.listen((e){
    headerDiv.style.height = (window.innerWidth.toInt()/3).toString() + 'px';
    titleImage.style.width = (window.innerWidth.toInt()/3).toString() + 'px';
    titleImage.style.paddingTop = (window.innerWidth.toInt()/24).toString() + 'px';
  });

  window.onScroll.listen((e) {
    Element topNavbar = querySelector('#top-navbar');
    if(window.scrollY >= 100)
    {
      topNavbar.style.backgroundColor = 'rgba(200, 180, 135, 0.6)';
    } else {
      topNavbar.style.backgroundColor = 'rgba(200, 180, 135, 1)';
    }
    if(((window.innerWidth.toInt()/24) + (window.scrollY)/2) < headerDiv.clientHeight/2) {
      headerDiv.style.backgroundPositionY = (window.scrollY / 2).toString() + 'px';
      titleImage.style.paddingTop = ((window.innerWidth.toInt() / 24) + (window.scrollY) / 2).toString() + 'px';
    }
  });
}

// Adds modal content to modal div
addModal(String modalLinkId, String modalId, String page) async {
  Element modalDiv = querySelector('#modal');
  Element modalLink = querySelector('#' + modalLinkId);

  modalLink.attributes['data-toggle'] = 'modal';
  modalLink.attributes['data-target'] = '#' + modalId;
  String modalHtml = await HttpRequest.getString(page);
  await modalDiv.appendHtml(modalHtml, treeSanitizer: NodeTreeSanitizer.trusted);
  if (modalId == 'signinModal') {
    setSigninListeners();
  }
  if (modalId == 'signupModal') {
    setSignupListeners();
  }
}

// Add listeners to show the signin modal
setSigninListeners(){
  // Toggle forgot password content
  var forgotPasswordListener = querySelector('#forgotPassword');
  var forgotPasswordContent = querySelector('#forgotPasswordContent');
  var signinContent = querySelector('#signinContent');
  forgotPasswordListener.onClick.listen((e){
    forgotPasswordContent.style.visibility = 'visible';
    signinContent.style.visibility = 'hidden';
  });

  var forgotPwdCloseListener = querySelector('#close-forgot-password');
  forgotPwdCloseListener.onClick.listen((e){
    var forgotPasswordError = querySelector('#forgotPasswordError');
    forgotPasswordListener.innerHtml = '';
    forgotPasswordListener.style.display = 'none';
    forgotPasswordContent.style.visibility = 'hidden';
    signinContent.style.visibility = 'visible';
  });

  // Clear warning and error dialogs when X is clicked
  var signinCloseListener = querySelector('#close-signin');
  signinCloseListener.onClick.listen((e){
    var signinErr = querySelector('#signinError');
    var signinWarn = querySelector('#signinWarn');
    signinErr.innerHtml = '';
    signinErr.style.display = 'none';
    signinWarn.innerHtml = '';
    signinWarn.style.display = 'none';
  });
}

// Add listeners to the signup modal
setSignupListeners() {
  // Clear warning and error dialogs when X is clicked
  var signupCloseListener = querySelector('#close-signup');
  signupCloseListener.onClick.listen((e){
    var signupErr = querySelector('#signupError');
    signupErr.innerHtml = '';
    signupErr.style.display = 'none';
  });
}


