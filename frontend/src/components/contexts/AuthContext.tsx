import * as React from "react";
import { createContext, useContext, useState, useEffect } from "react";
import {
  AuthService,
  AuthResult,
} from "../../libs/firebase/AuthService";
import LoadingSpinner from "../parts/LoadingSpinner";

type AuthState = {
  isAuthorized: boolean;
  isAdmin: boolean;
  uid: string;
};

type ContextType = {
  state: AuthState;
  loading: boolean;
  signIn: (email: string, password: string) => Promise<AuthResult>;
  signOut: () => Promise<void>;
};

const service = new AuthService();
const initialState = { isAuthorized: false, isAdmin: false,  uid: "" };

const AuthContext = createContext({} as ContextType);

export function useAuth(): ContextType {
  return useContext(AuthContext);
}

export const AdminAuthProvider = ({ children }: any) => {
  const [state, setState] = useState(initialState);
  const [initializing, setInitializing] = useState(true);
  const [loading, setLoading] = useState(false);

  async function signIn(
    email: string,
    password: string
  ): Promise<AuthResult> {
    setLoading(true);
    const result = await service.signIn(email, password);
    setState({ isAuthorized: result.isSuccessful, isAdmin: result.isAdmin, uid: result.uid });
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
      setState({ isAuthorized: result.isSuccessful, isAdmin: result.isAdmin, uid: result.uid });
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
    <AuthContext.Provider value={values}>
      {!initializing && children}
    </AuthContext.Provider>
  );
};
