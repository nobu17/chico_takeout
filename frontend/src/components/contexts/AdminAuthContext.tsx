import * as React from "react";
import { createContext, useContext, useState, useEffect } from "react";
import {
  AdminAuthService,
  AdminAuthResult,
} from "../../libs/firebase/AdminAuthService";
import LoadingSpinner from "../../components/parts/LoadingSpinner";

type AdminAuthState = {
  isAuthorized: boolean;
  uid: string;
};

type ContextType = {
  state: AdminAuthState;
  loading: boolean;
  signIn: (email: string, password: string) => Promise<AdminAuthResult>;
  signOut: () => Promise<void>;
};

const service = new AdminAuthService();
const initialState = { isAuthorized: false, uid: "" };

const AdminAuthContext = createContext({} as ContextType);

export function useAdminAuth(): ContextType {
  return useContext(AdminAuthContext);
}

export const AdminAuthProvider = ({ children }: any) => {
  const [state, setState] = useState(initialState);
  const [initializing, setInitializing] = useState(true);
  const [loading, setLoading] = useState(false);

  async function signIn(
    email: string,
    password: string
  ): Promise<AdminAuthResult> {
    setLoading(true);
    const result = await service.signIn(email, password);
    setState({ isAuthorized: result.isSuccessful, uid: result.uid });
    setLoading(false);
    return result;
  }

  async function signOut(): Promise<void> {
    try {
      await service.signOut();
      setState({ ...state, isAuthorized: false });
    } catch (err) {
      console.error("failed to logout");
    }
  }

  useEffect(() => {
    service.onAuthStateChange((result) => {
      setInitializing(false);
      setState({ isAuthorized: result.isSuccessful, uid: result.uid });
    });
  }, []);

  const values = {
    state,
    loading,
    signIn,
    signOut,
  };

  if (initializing) {
    return <LoadingSpinner isLoading={true} />;
  }

  return (
    <AdminAuthContext.Provider value={values}>
      {!initializing && children}
    </AdminAuthContext.Provider>
  );
};
