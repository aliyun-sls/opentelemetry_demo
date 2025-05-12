// demoData.ts
export const ordersDemoData = {
  code: 200,
  message: "ok",
  data: [
    {
      id: 1,
      created_at: "2025-05-09T09:48:26.25+08:00",
      updated_at: "2025-05-09T09:48:26.25+08:00",
      deleted_at: null,
      user_id: 1,
      order_id: "8fdf9241-2be9-11f0-b7d1-8a30a50add1b",
      shipping_tracking_id: "bbcd4e93-40b3-442b-bbb7-3a6971fede29",
      currency_code: "USD",
      units: 189,
      nanos: 400000000,
      street_address: "1600 Amphitheatre Parkway",
      city: "Mountain View",
      state: "CA",
      country: "United States",
      zip_code: "94043",
      status: 1,
      logistics_status: 0,
      total_price: 0,
      order_details: [
        {
          order_detail_id: 1,
          product_id: "xxx",
          units: 111,
          nanos: 3333,
          user_id: 111,
          product_picture: "sss",
          product_name: "sss",
          description: "sss"
        },
        {
          order_detail_id: 2,
          product_id: "2xxx",
          units: 22222,
          nanos: 222,
          user_id: 111,
          product_picture: "222",
          product_name: "222",
          description: "222"
        }
      ]
    },
    {
      id: 2,
      created_at: "2025-05-09T09:48:26.25+08:00",
      updated_at: "2025-05-09T09:48:26.25+08:00",
      deleted_at: null,
      user_id: 1,
      order_id: "8fdf9241-2be9-11f0-b7d1-fdfdfff",
      shipping_tracking_id: "bbcd4e93-40b3-442b-bbb7-3a6971fede29",
      currency_code: "USD",
      units: 3333,
      nanos: 400000000,
      street_address: "1600 Amphitheatre Parkway",
      city: "Mountain View",
      state: "CA",
      country: "United States",
      zip_code: "94043",
      status: 1,
      logistics_status: 0,
      total_price: 0,
      order_details: [
        {
          order_detail_id: 1,
          product_id: "xxx",
          units: 111,
          nanos: 3333,
          user_id: 111,
          product_picture: "sss",
          product_name: "sss",
          description: "sss"
        }
      ]
    }
  ]
};