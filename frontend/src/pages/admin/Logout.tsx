import React from "react";
import { useNavigate } from "react-router-dom";
import { useAdminAuth } from "../../components/contexts/AdminAuthContext";
import { useEffect } from "react";

export default function AdminLogout() {
  const navigate = useNavigate();
  const { signOut } = useAdminAuth();

  useEffect(() => {
    const f = async () => {
      await signOut();
      navigate("/");
    };
    f();
  });

  return <></>;
}
