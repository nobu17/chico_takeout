import { UserInfo } from "../../hooks/UseUserInfo";

const storeKey = "userInfo";

const defaultUserInfo: UserInfo = {
  name: "",
  tel: "",
  email: "",
  memo: "",
};

export default class UserInfoStore {
  save(target: UserInfo) {
    try {
        const obj = JSON.parse(JSON.stringify(target));
        // remove memo not save target
        obj.memo = "";
        const jsonStr = JSON.stringify(obj)
        localStorage.setItem(storeKey, jsonStr);
    } catch (e: any) {
        console.error("failed set user", e);
      }
  }
  load(): UserInfo {
    try {
      const json = localStorage.getItem(storeKey);
      if (!json) {
        return defaultUserInfo;
      }
      return JSON.parse(json);
    } catch (e: any) {
      console.error("failed get user", e);
      return defaultUserInfo;
    }
  }
}
