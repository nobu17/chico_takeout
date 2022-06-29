import { useState } from "react";

export type UserInfo = {
  name: string;
  tel: string;
  email: string;
  memo: string;
};

const defaultUserInfo: UserInfo = {
  name: "",
  tel: "",
  email: "",
  memo: "",
};

export function useUserInfo() {
  const [userInfo, setUserInfo] = useState<UserInfo>(defaultUserInfo);
  const updateUserInfo = (request: UserInfo) => {
    setUserInfo({...userInfo, ...request});
  };

  return {
    userInfo,
    updateUserInfo,
  };
}
