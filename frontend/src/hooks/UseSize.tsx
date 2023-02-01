import { useTheme } from "@mui/material/styles";
import useMediaQuery from "@mui/material/useMediaQuery";

export const useSize = () => {
  const theme = useTheme();
  const isMobileSize = useMediaQuery(theme.breakpoints.down("sm"));
  return { isMobileSize };
};
