import {
    storage,
    storageRef,
    storageUploadBytes,
    storageGetDownloadURL,
  } from "./Firebase";
  import {
    getCompressImageFileAsync,
    getCompressThumbNailImageFileAsync,
  } from "../util/ImageUtil";
  
  export const uploadFile = async (dir: string, name: string, file: File) => {
    try {
      let refs = storageRef(storage, dir);
  
      // need to cache
      const metadata = {
        cacheControl: "public,max-age=100000",
        contentType: "image/jpeg",
      };
      refs = storageRef(refs, name);
      const snapshot = await storageUploadBytes(refs, file, metadata);
      return await storageGetDownloadURL(snapshot.ref);
    } catch (err) {
      console.error("upload is fail", err);
      throw err;
    }
  };
    
  export const uploadImageWithCompress = async (
    dir: string,
    name: string,
    file: File
  ) => {
    const newFile = await getCompressImageFileAsync(file);
    return await uploadFile(dir, name, newFile);
  };
  
  export const uploadImageWithCompressAndThumbNail = async (
    dir: string,
    name: string,
    file: File
  ) => {
    const ext = ".jpg";
    const newFile = await getCompressImageFileAsync(file);
    const thumbFile = await getCompressThumbNailImageFileAsync(file);
  
    const thumbName = name.replace(ext, "_thumb.jpg");
  
    const fileUrl = await uploadFile(dir, name, newFile);
    const thumbUrl = await uploadFile(dir, thumbName, thumbFile);
  
    return [fileUrl, thumbUrl];
  };