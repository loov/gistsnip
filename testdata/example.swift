class ManagerAssembly: Assembly {
    //gistsnip:start:gist1
    func assemble(container: Container) {
        //gistsnip:start:gist2
        //gistsnip:start:gist3
        container.autoregister(ChildManagering.self, initializer: ChildManager.init)
        //gistsnip:end:gist2
            .inObjectScope(.container)
        //gistsnip:end:gist1
        container.autoregister(UserManagering.self, initializer: UserManager.init)
            .inObjectScope(.container)
    }
    //gistsnip:end:gist3
}