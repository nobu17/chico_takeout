import { useState } from "react";
import { useAuth } from "../components/contexts/AuthContext";
import UserInfoStore from "../libs/apis/userInfo";

export type UserInfo = {
  name: string;
  tel: string;
  email: string;
  memo: string;
};

const store = new UserInfoStore();

const defaultUserInfo = (authEmail : string): UserInfo => {
  let data = store.load();
  if (authEmail !== "") {
    data.email = authEmail;
  }
  return data;
} 

export function useUserInfo() {
  const { state } = useAuth();
  const [userInfo, setUserInfo] = useState<UserInfo>(defaultUserInfo(state.email));
  const updateUserInfo = (request: UserInfo) => {
    setUserInfo({...userInfo, ...request});
    store.save(request);
  };

  return {
    userInfo,
    updateUserInfo,
  };
}
