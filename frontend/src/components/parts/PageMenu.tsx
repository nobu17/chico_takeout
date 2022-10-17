import { NavLink as RouterLink } from "react-router-dom";

import ListSubheader from "@mui/material/ListSubheader";
import List from "@mui/material/List";
import ListItemButton from "@mui/material/ListItemButton";
import ListItemIcon from "@mui/material/ListItemIcon";
import ListItemText from "@mui/material/ListItemText";
import Stack from "@mui/material/Stack";
import DoubleArrowIcon from "@mui/icons-material/DoubleArrow";
import CoffeeIcon from "@mui/icons-material/Coffee";
import AppleIcon from "@mui/icons-material/Apple";
import ScienceIcon from '@mui/icons-material/Science';
import CalendarMonthIcon from "@mui/icons-material/CalendarMonth";
import { Typography } from "@mui/material";

export type PageMenuProps = {
  title: string;
  icon: Icon;
  pageInfos: PageInfo[];
};

export type PageInfo = {
  title: string;
  url: string;
  disabled?: boolean;
};

export type Icon = "coffee" | "apple" | "calendar" | "science"

const getIcon = (input: Icon) => {
  switch (input) {
    case "coffee": {
      return <CoffeeIcon />;
    }
    case "apple": {
      return <AppleIcon />;
    }
    case "calendar": {
      return <CalendarMonthIcon />;
    }
    case "science": {
      return <ScienceIcon />;
    } 
    default: {
      return <CoffeeIcon />;
    }
  }
};

export default function PageMenu(props: PageMenuProps) {
  return (
    <>
      <List
        sx={{
          width: "100%",
          maxWidth: 460,
        }}
        component="nav"
        subheader={
          <ListSubheader
            component="div"
            sx={{
              color: "text.primary",
              fontSize: 20,
              textAlign: "center",
              fontWeight: "bold",
              pb: 2,
            }}
          >
            <Stack
              direction="row"
              justifyContent="center"
              alignItems="center"
              gap={1}
            >
              {getIcon(props.icon)}
              <Typography variant="h5">{props.title}</Typography>
            </Stack>
          </ListSubheader>
        }
      >
        {props.pageInfos.map((item, index) =>
          item.disabled ? (
            <ListItemButton key={index} disabled={true}>
              <ListItemIcon>
                <DoubleArrowIcon />
              </ListItemIcon>
              <ListItemText primary={item.title} />
            </ListItemButton>
          ) : (
            <ListItemButton key={index} component={RouterLink} to={item.url}>
              <ListItemIcon>
                <DoubleArrowIcon />
              </ListItemIcon>
              <ListItemText primary={item.title} />
            </ListItemButton>
          )
        )}
      </List>
    </>
  );
}
