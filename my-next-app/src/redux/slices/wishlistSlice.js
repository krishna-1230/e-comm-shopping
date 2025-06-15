import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import { wishlistAPI } from '../../services/api';

const initialState = {
  items: [],
  loading: false,
  error: null,
};

// Async thunks for API calls
export const fetchWishlist = createAsyncThunk(
  'wishlist/fetchWishlist',
  async (_, { rejectWithValue }) => {
    try {
      const response = await wishlistAPI.getWishlist();
      return response.data;
    } catch (error) {
      return rejectWithValue(error.response?.data?.message || 'Failed to fetch wishlist');
    }
  }
);

export const addItemToWishlist = createAsyncThunk(
  'wishlist/addItem',
  async (item, { rejectWithValue }) => {
    try {
      const response = await wishlistAPI.addToWishlist(item);
      return response.data;
    } catch (error) {
      return rejectWithValue(error.response?.data?.message || 'Failed to add item to wishlist');
    }
  }
);

export const removeWishlistItem = createAsyncThunk(
  'wishlist/removeItem',
  async (id, { rejectWithValue }) => {
    try {
      await wishlistAPI.removeFromWishlist(id);
      return id;
    } catch (error) {
      return rejectWithValue(error.response?.data?.message || 'Failed to remove item from wishlist');
    }
  }
);

export const clearWishlistItems = createAsyncThunk(
  'wishlist/clearWishlist',
  async (_, { rejectWithValue }) => {
    try {
      await wishlistAPI.clearWishlist();
      return true;
    } catch (error) {
      return rejectWithValue(error.response?.data?.message || 'Failed to clear wishlist');
    }
  }
);

export const toggleWishlistProduct = createAsyncThunk(
  'wishlist/toggleItem',
  async (product, { getState, dispatch }) => {
    const { wishlist } = getState();
    const existingItem = wishlist.items.find(item => item.id === product.id);
    
    if (existingItem) {
      await dispatch(removeWishlistItem(product.id));
    } else {
      await dispatch(addItemToWishlist(product));
    }
    
    return product.id;
  }
);

export const wishlistSlice = createSlice({
  name: 'wishlist',
  initialState,
  reducers: {
    setLoading: (state, action) => {
      state.loading = action.payload;
    },
    
    setWishlist: (state, action) => {
      state.items = action.payload;
    },
    
    clearError: (state) => {
      state.error = null;
    },
  },
  extraReducers: (builder) => {
    builder
      // Fetch wishlist
      .addCase(fetchWishlist.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchWishlist.fulfilled, (state, action) => {
        state.loading = false;
        state.items = action.payload;
      })
      .addCase(fetchWishlist.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload;
      })
      
      // Add to wishlist
      .addCase(addItemToWishlist.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(addItemToWishlist.fulfilled, (state, action) => {
        state.loading = false;
        state.items = action.payload;
      })
      .addCase(addItemToWishlist.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload;
      })
      
      // Remove from wishlist
      .addCase(removeWishlistItem.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(removeWishlistItem.fulfilled, (state, action) => {
        state.loading = false;
        state.items = state.items.filter(item => item.id !== action.payload);
      })
      .addCase(removeWishlistItem.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload;
      })
      
      // Clear wishlist
      .addCase(clearWishlistItems.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(clearWishlistItems.fulfilled, (state) => {
        state.loading = false;
        state.items = [];
      })
      .addCase(clearWishlistItems.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload;
      });
  },
});

export const { setLoading, setWishlist, clearError } = wishlistSlice.actions;

// Utility selector
export const isItemInWishlist = (state, productId) => {
  return state.wishlist.items.some(item => item.id === productId);
};

export default wishlistSlice.reducer; 