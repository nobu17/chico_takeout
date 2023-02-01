import { Button, Badge } from "@mui/material";
import ShoppingCartIcon from "@mui/icons-material/ShoppingCart";

type FloatingCartIconButtonProps = {
  count: number;
  onClick: () => void;
};

export default function FloatingCartIconButton(
  props: FloatingCartIconButtonProps
) {
  return (
    <>
      <Button
        style={{
          background: "white",
          border: "solid",
          borderRadius: 50,
          minWidth: 50,
          width: 50,
          height: 50,
          position: "fixed",
          bottom: 50,
          right: 30,
        }}
        className="z-depth-1 p-2 d-flex justify-content-center align-items-center"
        onClick={props.onClick}
      >
        <Badge color="secondary" badgeContent={props.count}>
          <ShoppingCartIcon style={{ fontSize: 28 }} className="text-primary" />
        </Badge>
      </Button>
    </>
  );
}
