import * as React from "react";
import { createContext, useContext, useState, useEffect } from "react";
import { AuthService, AuthResult } from "../../libs/firebase/AuthService";
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
  signInWithGoogle: () => Promise<void>;
  signInWithTwitter: () => Promise<void>;
  signOut: () => Promise<void>;
};

const service = new AuthService();
const initialState = { isAuthorized: false, isAdmin: false, uid: "" };

const AuthContext = createContext({} as ContextType);

export function useAuth(): ContextType {
  return useContext(AuthContext);
}

export const AdminAuthProvider = ({ children }: any) => {
  const [state, setState] = useState(initialState);
  const [initializing, setInitializing] = useState(true);
  const [loading, setLoading] = useState(false);

  async function signIn(email: string, password: string): Promise<AuthResult> {
    setLoading(true);
    const result = await service.signIn(email, password);
    setState({
      isAuthorized: result.isSuccessful,
      isAdmin: result.isAdmin,
      uid: result.uid,
    });
    setLoading(false);
    return result;
  }

  async function signInWithGoogle(): Promise<void> {
    setLoading(true);
    await service.signInWithGoogle();
    setLoading(false);
  }

  async function signInWithTwitter(): Promise<void> {
    setLoading(true);
    await service.signInWithTwitter();
    setLoading(false);
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
    service.getRedirectResult((result) => {
      setState({
        isAuthorized: result.isSuccessful,
        isAdmin: result.isAdmin,
        uid: result.uid,
      });
    })
    service.onAuthStateChange((result) => {
      setInitializing(false);
      setState({
        isAuthorized: result.isSuccessful,
        isAdmin: result.isAdmin,
        uid: result.uid,
      });
    });
  }, []);

  const values = {
    state,
    loading,
    signIn,
    signOut,
    signInWithGoogle,
    signInWithTwitter,
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
