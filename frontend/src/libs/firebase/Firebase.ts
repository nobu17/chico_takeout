import { initializeApp } from "firebase/app";
import { getAuth, signInWithEmailAndPassword, EmailAuthProvider, sendEmailVerification, getAdditionalUserInfo } from 'firebase/auth';
import { getStorage, ref, uploadBytes, getDownloadURL, updateMetadata, listAll } from "firebase/storage";

const firebaseConfig = {
  apiKey: process.env.REACT_APP_APIKEY,
  authDomain: process.env.REACT_APP_AUTHDOMAIN,
  databaseURL: process.env.REACT_APP_DATABASEURL,
  projectId: process.env.REACT_APP_PROJECT_ID,
  storageBucket: process.env.REACT_APP_STORAGE_BUCKET,
  messagingSenderId: process.env.REACT_APP_MESSAGING_SENDER_ID,
  appId: process.env.REACT_APP_APP_ID,
  measurementId: process.env.REACT_APP_MEASUREMENT_ID,
};
const firebase = initializeApp(firebaseConfig);

const auth_obj = getAuth();
const storage_obj = getStorage();


export default firebase;
export const auth = auth_obj;
export const sendMail = sendEmailVerification;
export const getAdditionalUser = getAdditionalUserInfo;
export const signIn = signInWithEmailAndPassword;
export const storage = storage_obj;
export const storageRef = ref;
export const storageUploadBytes = uploadBytes;
export const storageGetDownloadURL = getDownloadURL;
export const storageUpdateMetadata = updateMetadata;
export const storageListAll = listAll;
