library photoUpload;

import 'dart:html';

addPhotoListener(String panelId) {
  FileUploadInputElement addPhotoListener = querySelector("#${panelId} #trophy_upload");
  UListElement photoList = querySelector('#${panelId} #photo_list');
  InputElement photoName = querySelector('#${panelId} #photo_name');
  photoName.onFocus.listen((e){ photoName.value = '';});
  addPhotoListener.onChange.listen((e){
    var newFile = new LIElement();
    newFile.attributes['datavalue'] = addPhotoListener.files.last.name;
    newFile.text = photoName.value;
    newFile.attributes["id"] = "added_photo_${addPhotoListener.files.length}";
    photoList.children.add(newFile);
    photoName.value = "Enter photo name";

  });
}