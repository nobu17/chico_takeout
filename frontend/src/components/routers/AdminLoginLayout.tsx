import React from "react";
import { Outlet, Navigate } from "react-router-dom";
import { useAdminAuth } from "../contexts/AdminAuthContext";

export default function AdminLoginLayout() {
  const { state } = useAdminAuth();

  return state.isAuthorized ? <Navigate to="/admin" /> : <Outlet />;
}
