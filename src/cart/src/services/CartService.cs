// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0
using System.Diagnostics;
using System.Threading.Tasks;
using System;
using Grpc.Core;
using cart.cartstore;
using OpenFeature;
using Oteldemo;
using Microsoft.Extensions.Logging;

namespace cart.services;

public class CartService : Oteldemo.CartService.CartServiceBase
{
    private static readonly Empty Empty = new();
    private readonly Random random = new Random();
    private readonly ICartStore _badCartStore;
    private readonly ICartStore _cartStore;
    private readonly IFeatureClient _featureFlagHelper;
    private readonly ILogger<CartService> _logger;

    public CartService(ICartStore cartStore, ICartStore badCartStore, IFeatureClient featureFlagService, ILogger<CartService> logger)
    {
        _badCartStore = badCartStore;
        _cartStore = cartStore;
        _featureFlagHelper = featureFlagService;
        _logger = logger;
    }

    public override async Task<Empty> AddItem(AddItemRequest request, ServerCallContext context)
    {
        var activity = Activity.Current;
        activity?.SetTag("app.user.id", request.UserId);
        activity?.SetTag("app.product.id", request.Item.ProductId);
        activity?.SetTag("app.product.quantity", request.Item.Quantity);

        _logger.LogInformation("用户 {UserId} 正在添加商品到购物车: 商品ID {ProductId}, 数量 {Quantity}", 
            request.UserId, request.Item.ProductId, request.Item.Quantity);

        try
        {
            await _cartStore.AddItemAsync(request.UserId, request.Item.ProductId, request.Item.Quantity);
            _logger.LogInformation("用户 {UserId} 成功添加商品到购物车: 商品ID {ProductId}, 数量 {Quantity}", 
                request.UserId, request.Item.ProductId, request.Item.Quantity);

            return Empty;
        }
        catch (RpcException ex)
        {
            _logger.LogError(ex, "用户 {UserId} 添加商品到购物车失败: 商品ID {ProductId}, 数量 {Quantity}", 
                request.UserId, request.Item.ProductId, request.Item.Quantity);
            activity?.AddException(ex);
            activity?.SetStatus(ActivityStatusCode.Error, ex.Message);
            throw;
        }
    }

    public override async Task<Cart> GetCart(GetCartRequest request, ServerCallContext context)
    {
        var activity = Activity.Current;
        activity?.SetTag("app.user.id", request.UserId);
        activity?.AddEvent(new("Fetch cart"));

        _logger.LogInformation("用户 {UserId} 正在获取购物车内容", request.UserId);

        try
        {
            var cart = await _cartStore.GetCartAsync(request.UserId);
            var totalCart = 0;
            foreach (var item in cart.Items)
            {
                totalCart += item.Quantity;
            }
            activity?.SetTag("app.cart.items.count", totalCart);

            _logger.LogInformation("用户 {UserId} 成功获取购物车内容: 商品总数 {TotalItems}", 
                request.UserId, totalCart);

            return cart;
        }
        catch (RpcException ex)
        {
            _logger.LogError(ex, "用户 {UserId} 获取购物车内容失败", request.UserId);
            activity?.AddException(ex);
            activity?.SetStatus(ActivityStatusCode.Error, ex.Message);
            throw;
        }
    }

    public override async Task<Empty> EmptyCart(EmptyCartRequest request, ServerCallContext context)
    {
        var activity = Activity.Current;
        activity?.SetTag("app.user.id", request.UserId);
        activity?.AddEvent(new("Empty cart"));

        _logger.LogInformation("用户 {UserId} 正在清空购物车", request.UserId);

        try
        {
            if (await _featureFlagHelper.GetBooleanValueAsync("cartFailure", false))
            {
                _logger.LogWarning("用户 {UserId} 触发了购物车故障注入", request.UserId);
                await _badCartStore.EmptyCartAsync(request.UserId);
            }
            else
            {
                await _cartStore.EmptyCartAsync(request.UserId);
            }

            _logger.LogInformation("用户 {UserId} 成功清空购物车", request.UserId);
        }
        catch (RpcException ex)
        {
            _logger.LogError(ex, "用户 {UserId} 清空购物车失败", request.UserId);
            Activity.Current?.AddException(ex);
            Activity.Current?.SetStatus(ActivityStatusCode.Error, ex.Message);
            throw;
        }

        return Empty;
    }
}
