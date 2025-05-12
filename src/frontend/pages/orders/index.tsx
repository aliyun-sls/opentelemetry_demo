
import { useState, useEffect } from 'react';
import Link from 'next/link';
import styles from '../../styles/orders.module.css';
import { ordersDemoData, Order } from '../../utils/demoData';
import Header from '../../components/Header/Header';

const OrdersPage = () => {
  const [orders, setOrders] = useState<Order[]>([]);

  useEffect(() => {
    fetchOrders();
  }, []);

  const fetchOrders = async () => {
    const useDemoData = true;
    if (useDemoData) {
      setOrders(ordersDemoData.data);
      return;
    }

    try {
      const response = await fetch('http://order:8080/order/list');

      if (!response.ok) {
        setOrders([]);
      }

      const contentType = response.headers.get("content-type");
      if (!contentType || !contentType.includes("application/json")) {
        setOrders([]);
      }

      const data = await response.json();
      setOrders(data);
    } catch (error) {
      console.error('Failed to fetch orders:', error);
      setOrders([]);
    }
  };

  return (
    <>
      <Header />
      <div className={styles.container}>
        <h1 className={styles.title}>我的订单</h1>
        {orders.length === 0 ? (
          <p className={styles.emptyMessage}>暂无订单数据，请稍后再试或联系管理员。</p>
        ) : (
          <ul className={styles.ordersList}>
            {orders.map((order) => (
              <li key={order.id} className={styles.orderItem}>
                <Link href={`/orders/${order.id}`} legacyBehavior>
                  <a className={styles.orderLink}>
                    <div>订单 #{order.id}</div>
                    <div>总金额: {(order.units + order.nanos / 1e9).toFixed(2)} {order.currency_code}</div>
                    <div>商品数量: {order.order_details.length}</div>
                    <div>状态: {order.status === 1 ? '已完成' : '进行中'}</div>
                  </a>
                </Link>
              </li>
            ))}
          </ul>
        )}
      </div>
    </>
  );
};

export default OrdersPage;