import React from "react";
import { Navigate, Outlet } from "react-router-dom";
import { useAdminAuth } from "../contexts/AdminAuthContext";

export default function AdminAuthLayout() {
  const { state } = useAdminAuth();

  return state.isAuthorized ? <Outlet /> : <Navigate to="/admin/sign_in" />;
}
