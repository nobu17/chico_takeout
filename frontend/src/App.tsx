import React from "react";
import { Route, Routes, BrowserRouter } from "react-router-dom";
import {
  createTheme,
  responsiveFontSizes,
  ThemeProvider,
} from "@mui/material/styles";

import Container from "@mui/material/Container";
import Header from "./components/layouts/Header";
import AdminHome from "./pages/admin/Home";
import ItemKind from "./pages/admin/items/ItemKind";
import "./App.css";

let theme = createTheme();
theme = responsiveFontSizes(theme);

function App() {
  return (
    <ThemeProvider theme={theme}>
      <BrowserRouter>
        <Header />
        <Container maxWidth="lg">
          <Routes>
            <Route path="/admin" element={<AdminHome />} />
            <Route path="/admin/items/kind" element={<ItemKind />} /> 
          </Routes>
        </Container>
      </BrowserRouter>
    </ThemeProvider>
  );
}

export default App;
