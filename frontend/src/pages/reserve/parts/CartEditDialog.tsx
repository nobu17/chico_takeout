import * as React from "react";
import {
  Typography,
  AppBar,
  Toolbar,
  Slide,
  Dialog,
  Grid,
  Paper,
  Stack,
  Button,
  Select,
  MenuItem,
  DialogActions,
  DialogContent,
} from "@mui/material";
import { SelectChangeEvent } from "@mui/material/Select";
import { TransitionProps } from "@mui/material/transitions";
import {
  Cart,
  ItemInfo,
  OptionItemInfo,
  ItemRequest,
  useItemCart,
} from "../../../hooks/UseItemCart";
import { useSize } from "../../../hooks/UseSize";
import {
  getSubTotalPriceFromReq,
  getOptionTotalPriceFromReq,
  getTotalPriceFromCart,
} from "../../../libs/util/ItemCalc";
import OptionItemSelectDialog from "./OptionItemSelectDialog";

const Transition = React.forwardRef(function Transition(
  props: TransitionProps & {
    children: React.ReactElement;
  },
  ref: React.Ref<unknown>
) {
  return <Slide direction="up" ref={ref} {...props} />;
});

type CartEditDialogProps = {
  allItems: CategoryItems[];
  cart: Cart;
  open: boolean;
  onCancel: () => void;
  onSubmit: (cart: Cart) => void;
};

type CategoryItems = {
  title: string;
  items: ItemInfo[];
};

export default function CartEditDialog(props: CartEditDialogProps) {
  const { isMobileSize } = useSize();
  // copy item as draft (set init value is only attempt first time)
  const { cart, updateCartAsDraft, resetCart, getActivatedCartFromDraft } =
    useItemCart(props.cart);

  // if parent cart is updated, hook value is needed update
  React.useEffect(() => {
    resetCart(props.cart); // will deep copy is created as draft
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [props.cart]);

  const handleUpdate = (item: ItemRequest) => {
    updateCartAsDraft(item);
  };
  const handleSubmit = () => {
    // remove 0 quantity item
    const activated = getActivatedCartFromDraft();
    props.onSubmit(activated);
  };
  const handleCancel = () => {
    // discard edit items
    resetCart(props.cart);
    props.onCancel();
  };
  const handleReset = () => {
    // clear all items
    resetCart();
  };
  return (
    <>
      <Dialog
        open={props.open}
        fullWidth={!isMobileSize}
        maxWidth="md"
        fullScreen={isMobileSize}
        TransitionComponent={Transition}
      >
        <AppBar sx={{ position: "relative", mb: 1 }}>
          <Toolbar>
            <Typography sx={{ ml: 2, flex: 1 }} variant="h6" component="div">
              カート情報
            </Typography>
            <Button variant="contained" color="secondary" onClick={handleReset}>
              空にする
            </Button>
          </Toolbar>
        </AppBar>
        <DialogContent dividers={true} sx={{ px: 1, py: 0, m: 0 }}>
          <Stack spacing={1}>
            {props.open ? (
              Object.keys(cart.items).map((key) => {
                const req = cart.items[key];
                return (
                  <CartItemEditCard
                    key={key}
                    item={req}
                    onUpdate={handleUpdate}
                  ></CartItemEditCard>
                );
              })
            ) : (
              <></>
            )}
          </Stack>
        </DialogContent>
        <DialogActions>
          {props.open ? (
            <CartItemSummaryCard cart={cart}></CartItemSummaryCard>
          ) : (
            <></>
          )}
          <Button
            variant="contained"
            color="primary"
            sx={{ mr: 2, ml: 4 }}
            onClick={handleSubmit}
          >
            確定
          </Button>
          <Button variant="contained" color="error" onClick={handleCancel}>
            キャンセル
          </Button>
        </DialogActions>
      </Dialog>
    </>
  );
}

type CartItemEditCardProps = {
  item: ItemRequest;
  onUpdate: (item: ItemRequest) => void;
};

function CartItemEditCard(props: CartItemEditCardProps) {
  const req = props.item;
  const { isMobileSize } = useSize();
  const [open, setOpen] = React.useState(false);

  const handleQuantityChange = (event: SelectChangeEvent<number>) => {
    const quantity = event.target.value as number;
    if (quantity >= 0) {
      props.onUpdate({
        item: props.item.item,
        quantity: quantity,
        selectOptions: props.item.selectOptions,
      });
    }
  };

  const handleOptionClick = () => {
    setOpen(true);
  };

  const handleOptionClose = (options: OptionItemInfo[]) => {
    setOpen(false);
    props.onUpdate({
      item: props.item.item,
      quantity: props.item.quantity,
      selectOptions: options,
    });
  };

  return (
    <Paper
      elevation={0}
      sx={{
        p: 1,
        flexGrow: 1,
        px: isMobileSize ? 1 : 3,
        borderColor: "grey.500",
        borderBottom: 1,
        backgroundColor: (theme) =>
          theme.palette.mode === "dark" ? "#1A2027" : "#fff",
      }}
    >
      <Grid container spacing={2}>
        <Grid item xs={12} sm container>
          <Grid item xs container direction="column" spacing={2}>
            <Grid item xs>
              <Typography
                gutterBottom
                variant="subtitle2"
                component="div"
                sx={{ mr: 2 }}
              >
                {req.item.name}
              </Typography>
              <Stack direction="row" spacing={1}>
                <Typography variant="subtitle2" sx={{ mt: 2 }}>
                  注文数:
                </Typography>
                <Select
                  size="small"
                  value={req.quantity}
                  label="個数"
                  sx={{
                    minWidth: 50,
                    ml: 2,
                  }}
                  onChange={handleQuantityChange}
                >
                  {[...Array(req.item.max + 1)].map((_, i) => (
                    <MenuItem key={i} value={i}>
                      {i}
                    </MenuItem>
                  ))}
                </Select>
                {req.item.options.length > 0 ? (
                  <Button
                    color="error"
                    variant="contained"
                    size="small"
                    onClick={handleOptionClick}
                  >
                    オプション
                    {req.selectOptions.length > 0
                      ? "(" + req.selectOptions.length + ")"
                      : ""}
                  </Button>
                ) : (
                  <></>
                )}
              </Stack>
            </Grid>
          </Grid>
          <Grid item>
            <Stack
              spacing={2}
              direction="column"
              justifyContent="flex-end"
              alignItems="flex-end"
            >
              <Typography variant="subtitle1" component="div" align="right">
                ¥ {req.item.price}{" "}
                {req.selectOptions.length > 0
                  ? "(+¥ " + getOptionTotalPriceFromReq(req) + ")"
                  : ""}
              </Typography>
              <Typography variant="subtitle1" component="div" align="right">
                小計 ¥ {getSubTotalPriceFromReq(req)}
              </Typography>
            </Stack>
          </Grid>
        </Grid>
      </Grid>
      <OptionItemSelectDialog
        open={open}
        items={req.item.options}
        selected={req.selectOptions}
        onClose={handleOptionClose}
      ></OptionItemSelectDialog>
    </Paper>
  );
}

type CartItemSummaryCardProps = {
  cart: Cart;
};

function CartItemSummaryCard(props: CartItemSummaryCardProps) {
  const { isMobileSize } = useSize();
  return (
    <Paper
      elevation={0}
      sx={{
        p: 1,
        flexGrow: 1,
        px: isMobileSize ? 1 : 3,
      }}
    >
      <Typography variant="h5" component="div" align="right">
        合計 ¥ {getTotalPriceFromCart(props.cart)}
      </Typography>
    </Paper>
  );
}
