import React from "react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../../components/contexts/AuthContext";
import { useEffect } from "react";

export default function Logout() {
  const navigate = useNavigate();
  const { signOut } = useAuth();

  useEffect(() => {
    const f = async () => {
      await signOut();
      navigate("/");
    };
    f();
  });

  return <></>;
}
