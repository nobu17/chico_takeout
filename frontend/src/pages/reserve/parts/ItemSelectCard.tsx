import * as React from "react";
import {
  Typography,
  AppBar,
  Toolbar,
  Slide,
  DialogActions,
  DialogContent,
  Button,
  Stack,
  Box,
  Dialog,
  Card,
  CardActions,
  CardActionArea,
  CardContent,
  CardMedia,
} from "@mui/material";
import Badge, { BadgeProps } from "@mui/material/Badge";
import { styled } from "@mui/material/styles";

import ShoppingCartIcon from "@mui/icons-material/ShoppingCart";
import { TransitionProps } from "@mui/material/transitions";
import LineBreak from "../../../components/parts/LineBreak";
import Counter from "../../../components/parts/Counter";
import OptionItemSelectDialog from "./OptionItemSelectDialog";
import {
  ItemInfo,
  ItemRequest,
  OptionItemInfo,
} from "../../../hooks/UseItemCart";
import { getImageUrl } from "../../../libs/util/ImageUtil";
import { useSize } from "../../../hooks/UseSize";
import { getOptionTotalPrice } from "../../../libs/util/ItemCalc";

type ItemSelectCardProps = {
  item: ItemInfo;
  quantity: number;
  selectedOptions: OptionItemInfo[];
  onChanged: callback;
};
interface callback {
  (item: ItemRequest): void;
}

export default function ItemSelectCard(props: ItemSelectCardProps) {
  const [open, setOpen] = React.useState(false);

  const handleEditStart = () => {
    setOpen(true);
  };

  const handleSubmit = (item: ItemRequest) => {
    setOpen(false);
    props.onChanged(item);
  };

  const handleCancel = () => {
    setOpen(false);
  };

  return (
    <>
      <Card sx={{ maxWidth: 400, height: "100%" }}>
        <CardActionArea sx={{ height: "100%" }} onClick={handleEditStart}>
          <CardMedia
            component="img"
            height="140"
            sx={{ objectFit: "contain" }}
            image={getImageUrl(props.item.imageUrl)}
          />
          <CardContent
            sx={{
              height: "55px",
              display: "flex",
              flexDirection: "column",
              px: 1,
              pt: 1,
              pb: 0,
            }}
          >
            <Typography gutterBottom variant="body2" component="div">
              {props.item.name}
            </Typography>
          </CardContent>
          <CardActions disableSpacing>
            <Typography sx={{ ml: 1, mb: 0, mt: 0 }} variant="subtitle1">
              {" "}
              ¥ {props.item.price.toLocaleString()}
            </Typography>
            <CartBadge count={props.quantity}></CartBadge>
          </CardActions>
        </CardActionArea>
      </Card>
      <ItemSelectDialog
        item={props.item}
        quantity={props.quantity}
        selectedOptions={props.selectedOptions}
        open={open}
        onCancel={handleCancel}
        onSubmit={handleSubmit}
      ></ItemSelectDialog>
    </>
  );
}

const StyledBadge = styled(Badge)<BadgeProps>(({ theme }) => ({
  "& .MuiBadge-badge": {
    right: 0,
    top: 13,
    border: `2px solid ${theme.palette.background.paper}`,
    padding: "0 4px",
  },
}));

type CartBadgeProps = {
  count: number;
};

function CartBadge(props: CartBadgeProps) {
  if (props.count <= 0) {
    return <></>;
  }
  return (
    <StyledBadge
      badgeContent={props.count}
      color="secondary"
      sx={{ mr: 1, marginLeft: "auto", zIndex: 0, }}
    >
      <ShoppingCartIcon sx={{ mr: 0 }} />
    </StyledBadge>
  );
}

type ItemSelectDialogProps = {
  open: boolean;
  item: ItemInfo;
  quantity: number;
  selectedOptions: OptionItemInfo[];
  onSubmit: callback;
  onCancel: () => void;
};

const Transition = React.forwardRef(function Transition(
  props: TransitionProps & {
    children: React.ReactElement;
  },
  ref: React.Ref<unknown>
) {
  return <Slide direction="up" ref={ref} {...props} />;
});

function ItemSelectDialog(props: ItemSelectDialogProps) {
  const { isMobileSize } = useSize();
  const [optionOpen, setOptionOpen] = React.useState(false);
  const [options, setOptions] = React.useState(props.selectedOptions);
  const [quantity, setQuantity] = React.useState(props.quantity);

  // need update from cart button update
  React.useEffect(() => {
    setOptions(props.selectedOptions);
  }, [props.selectedOptions]);
  React.useEffect(() => {
    setQuantity(props.quantity);
  }, [props.quantity]);

  const onCountChanged = (count: number) => {
    setQuantity(count);
  };

  const handleOptionClick = () => {
    setOptionOpen(true);
  };

  const handleOptionClose = (selected: OptionItemInfo[]) => {
    setOptions(selected);
    setOptionOpen(false);
  };

  const handleSubmit = () => {
    props.onSubmit({
      item: props.item,
      quantity: quantity,
      selectOptions: options,
    });
  };

  const handleCancel = () => {
    props.onCancel();
  };

  return (
    <>
      <Dialog
        open={props.open}
        fullWidth={!isMobileSize}
        maxWidth="sm"
        fullScreen={isMobileSize}
        TransitionComponent={Transition}
        scroll="paper"
      >
        <AppBar sx={{ position: "relative" }}>
          <Toolbar>
            <Typography sx={{ ml: 2, flex: 1 }} variant="h6" component="div">
              商品情報
            </Typography>
          </Toolbar>
        </AppBar>
        <DialogContent dividers={true}>
          <Card sx={{ maxWidth: 900, mt: 2 }}>
            <CardMedia
              component="img"
              height="230"
              sx={{ objectFit: "contain" }}
              image={getImageUrl(props.item.imageUrl)}
            />
            <CardContent>
              <Typography gutterBottom variant="h6" component="div">
                {props.item.name}
              </Typography>
              <Typography variant="body2" color="text.secondary">
                <LineBreak msg={props.item.memo} />
              </Typography>
            </CardContent>
            <CardActions>
              <Stack>
                <Typography sx={{ ml: 1 }} variant="h5">
                  {" "}
                  ¥ {props.item.price.toLocaleString()}
                </Typography>
                {options.length < 1 ? (
                  <></>
                ) : (
                  <Typography sx={{ ml: 1 }} variant="subtitle2">
                    {" "}
                    {getOptionTotal(options)}
                  </Typography>
                )}
              </Stack>

              <Box style={{ marginLeft: "auto" }}>
                <Counter
                  count={quantity}
                  max={props.item.max}
                  onChanged={onCountChanged}
                ></Counter>
              </Box>
            </CardActions>
            {props.item.options.length <= 0 ? (
              <></>
            ) : (
              <Box sx={{ m: 2 }}>
                <Button
                  variant="contained"
                  color="error"
                  fullWidth
                  onClick={() => handleOptionClick()}
                >
                  オプション {getOptionCountString(options)}
                </Button>
              </Box>
            )}
            <OptionItemSelectDialog
              open={optionOpen}
              items={props.item.options}
              selected={options}
              onClose={handleOptionClose}
            />
          </Card>
        </DialogContent>

        <DialogActions>
          <Button
            variant="contained"
            color="primary"
            sx={{ mr: 1, width: 110 }}
            onClick={handleSubmit}
          >
            確定
          </Button>
          <Button variant="contained" sx={{ width: 110 }} color="error" onClick={handleCancel}>
            キャンセル
          </Button>
        </DialogActions>
      </Dialog>
    </>
  );
}

const getOptionCountString = (selected: OptionItemInfo[]): string => {
  if (selected.length < 1) {
    return "";
  }
  return `(${selected.length} 個)`;
};

const getOptionTotal = (options: OptionItemInfo[]): string => {
  const totalStr = getOptionTotalPrice(options).toLocaleString();
  return `  (+¥ ${totalStr})`;
};
