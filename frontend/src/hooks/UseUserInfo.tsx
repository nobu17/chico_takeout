import { useState } from "react";
import UserInfoStore from "../libs/apis/userInfo";

export type UserInfo = {
  name: string;
  tel: string;
  email: string;
  memo: string;
};

const store = new UserInfoStore();

export function useUserInfo() {
  const [userInfo, setUserInfo] = useState<UserInfo>(store.load());
  const updateUserInfo = (request: UserInfo) => {
    setUserInfo({...userInfo, ...request});
    store.save(request);
  };

  return {
    userInfo,
    updateUserInfo,
  };
}
